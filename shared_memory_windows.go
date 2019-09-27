package shm

import (
	"syscall"
	"unsafe"
)

//映射地址对应mapid
var AccessAddressesToHandles map[uintptr]syscall.Handle = map[uintptr]syscall.Handle{}

func CreateSharedMemory(size uint32,sharedMemId *uint16) uintptr{
	hMapFile,err := syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_EXECUTE_READWRITE, 0, size, sharedMemId)
	if err != nil {
		//util.WarnF("[SHDMEM] cannot create file mapping for smid:%v error %v", *sharedMemId, err)
		return uintptr(0)
	}

	accessAddress,_ := syscall.MapViewOfFile(hMapFile,
		syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ, 0, 0,0)
	AccessAddressesToHandles[accessAddress] = hMapFile
	return accessAddress
}

//下面注释有另外一种实现方法
func AccessSharedMemory(sharedMemId *uint16) uintptr {
	procOpenFileMapping := syscall.NewLazyDLL("kernel32.dll").NewProc("OpenFileMappingW")
	//执行syscall，出现exception，一般是参数类型错误导致
	r0,_,e1 := syscall.Syscall(procOpenFileMapping.Addr(), 3,
		uintptr(syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ), uintptr(0), uintptr(unsafe.Pointer(sharedMemId)))
	handle := syscall.Handle(r0)
	if e1 != 0 {
		//util.WarnF("[SHDMEM] OpenFileMappingW failed %v error:%v", *sharedMemId, e1)
		return uintptr(0)
	}

	accessAddress,_ := syscall.MapViewOfFile(handle,
		syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ, 0, 0,0)
	AccessAddressesToHandles[accessAddress] = handle
	return accessAddress
}

func CloseShardMemory(accessAddress uintptr) bool {
	 err := syscall.UnmapViewOfFile(accessAddress)
	 if err != nil {
		 //util.WarnF("[SHDMEM] UnmapViewOfFile failed err:%v", err)
	 	return false
	 }

	value, ok := AccessAddressesToHandles[accessAddress]
	if ok {
		err = syscall.CloseHandle(value)
		if err != nil {
			//util.WarnF("[SHDMEM] CloseHandle failed err:%v", err)
		}
		delete(AccessAddressesToHandles, accessAddress)
		return true
	}

	return false
}

//this method does nothing under windows, destroying is automatic.
func destroySharedMemory(sharedMemId *uint16, force bool) {
}

/*
func AccessSharedMemory(sharedName string) uintptr {
	shareNameByte,_ :=syscall.BytePtrFromString("Global/" + sharedName)
	var shareNameId = uintptr(unsafe.Pointer(shareNameByte))

	procOpenFileMapping := syscall.NewLazyDLL("kernel32.dll").NewProc("OpenFileMappingW")
	hMapFile,_,err := procOpenFileMapping.Call(syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ, 0, shareNameId)
	if err != nil {
		util.WarnF("[SHDMEM] OpenFileMappingW failed error:%v", err)
		return uintptr(0)
	}

	accessAddress,err1 := syscall.MapViewOfFile(syscall.Handle(hMapFile),
		syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ, 0, 0,0)
	if err1 != nil {
		util.WarnF("[SHDMEM] MapViewOfFile failed error:%v", err)
		return uintptr(0)
	}
	AccessAddressesToHandles[accessAddress] = syscall.Handle(hMapFile)
	return accessAddress

}
 */