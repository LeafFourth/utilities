package utilities

import "archive/zip"
import "errors"
import "io"
import "os"

func WriteFile(reader io.Reader, path string) error {
	mode := os.O_CREATE | os.O_WRONLY | os.O_TRUNC;

	dst, err := os.OpenFile(path, mode, 0755);
	if (err != nil) {
		return errors.New("create file error:" + path);
	}

	var result error = nil;
	for {
		var tmp = [1024]byte{}
		var bytes []byte = tmp[0:];
		n, err2 := reader.Read(bytes);
		if err2 != nil && err2 != io.EOF {
			result = errors.New("read error:" + path);
			break;
		}
		_, err3 := dst.Write(bytes);
		if (err3 != nil) {
			result = errors.New("read error:" + path);
			break;
		}
		if n == 0 {
			break;
		}
	}

	dst.Close();

	if result != nil {
		os.Remove(path);
	}
	
	return result;
}

func UnzipFile(zipFile string, dst string) error {
	if len(dst) <= 0 {
		return errors.New("empty path");
	}

	file, err := os.Open(dst);
	if err != nil {
		return errors.New("path not exists:" + zipFile);
	}
	defer file.Close();

	info, err2 := file.Stat();
	if err2 != nil {
		return errors.New("io error:" + zipFile);
	}

	if !info.IsDir() {
		return errors.New("not dir:" + zipFile);
	}

	if dst[len(dst) - 1] != '/' {
		dst += "/";
	}

	zipReader, err3 := zip.OpenReader(zipFile);
	if err3 != nil {
		return errors.New("zip not open:" + zipFile);
	}
	defer zipReader.Close();

	for _, f := range zipReader.File {
		subPath := dst + f.Name;
		if subPath[len(subPath) - 1] == '/' {
			err4 :=  os.MkdirAll(subPath, 0644);
			if err4 != nil {
				return errors.New("unzip error:" + zipFile);
			}
			continue;
		}


		reader, err5 := f.Open();
		if err5 != nil {
			return errors.New("unzip error:" + zipFile);
		}
		err6 := WriteFile(reader, subPath);
		reader.Close();
		if err6 != nil {
			return errors.New("unzip error:" + zipFile);
		}
	}
	return nil;
}