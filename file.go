package utilities

import "archive/zip"
import "errors"
import "io"
import "os"
import "path/filepath"

func WriteFile(reader io.Reader, path string) error {
	mode := os.O_CREATE | os.O_WRONLY | os.O_TRUNC;

	dst, err := os.OpenFile(path, mode, 0644);
	if (err != nil) {
		return errors.New("create file error:" + path);
	}

	var result error = nil;
	_, result = io.Copy(dst, reader);

	dst.Close();

	if result != nil &&  result != io.EOF {
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
		subPath := filepath.Join(dst, f.Name);	
		if f.Name[len(f.Name) - 1] == '/' {
			err4 :=  os.MkdirAll(subPath, 0644);
			if err4 != nil {
				return errors.New("unzip " + subPath + " create dir error:" + subPath);
			}
			continue;
		}

		parent := filepath.Dir(subPath);
		if err := os.MkdirAll(parent, 0644); err != nil {
			return err;
		}

		reader, err5 := f.Open();
		if err5 != nil {
			return errors.New("unzip open file error:" + subPath);
		}
		err6 := WriteFile(reader, subPath);
		reader.Close();
		if err6 != nil {
			return err6;
		}
	}
	return nil;
}