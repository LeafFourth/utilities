package utilities

import "fmt"
import "os"
import "unsafe"
import "sync"
import "sync/atomic"

type spinLock32 struct {
	val *[8]byte
	off uint
}

type LogWriter struct {
	mutex *sync.Mutex
	file *os.File
}

func alignMem(addr uintptr,base uint) uintptr {
	if (base != 4) {
		return 0;
	}

	ret := addr - addr % (1 << base);
	return ret;
}

func newSpinLock32() spinLock32 {
	ret := spinLock32{};
	ret.val = new([8]byte);
	if ret.val == nil {
		return ret;
	}
	fmt.Println(ret.val);
 	p := unsafe.Pointer(ret.val);
	p2 := alignMem(uintptr(p), 4) - uintptr(p)
	ret.off = uint(p2);
	return ret;
}

func (self spinLock32) lock() {
	addr := unsafe.Pointer(&self.val[self.off])
	addr2 := (*uint32)(addr);

	for ;atomic.CompareAndSwapUint32(addr2, 0, 1); {
	}
}

func (self spinLock32) unlock() {
	addr := unsafe.Pointer(&self.val[self.off])
	addr2 := (*int32)(addr);

	atomic.StoreInt32(addr2, 0)
}


func (self *LogWriter)Write(p []byte) (n int, err error) {
	self.mutex.Lock();
	defer self.mutex.Unlock();

	fp := self.file;

	if fp == nil {
		fmt.Println(p);
	} else {
		n, err := fp.Write(p);
		if err != nil {
		    fmt.Println("log file error:", fp.Name());
		}
		return n, err;
	}
	return len(p), nil;
}

func NewLogWriter() *LogWriter {
	ret := new(LogWriter);
	if ret == nil {
		return nil;
	}

	ret.mutex = new(sync.Mutex);
	if ret.mutex == nil {
		return nil;
	}

	return ret;
}

func (self *LogWriter)SetLogPath(path string) {
	self.mutex.Lock()
	if self.file != nil {
		self.file.Close()
		self.file = nil
	}
	if len(path) > 0 {
		fp, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDONLY , 0644)
		if err != nil {
			fmt.Println(err);
		} else {
			self.file = fp
		}

	}
	self.mutex.Unlock();
}


