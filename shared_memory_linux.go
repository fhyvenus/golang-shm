package shm

import (
	"syscall"
)

//映射地址对应mapid
var SharedMemIdsToShmids map[uint16]uintptr = map[uint16]uintptr{}

const (
	IPC_CREATE = 01000	//create entry if key does not exist
	IPC_EXCL = 02000	//fail if key exists
	IPC_NOWAIT = 04000	//error if request must wait

	IPC_RMID = 0 //remove identifier
	IPC_SET = 1 //set options
	IPC_STAT = 2 //get options
)

func CreateSharedMemory(size uint32,sharedMemId *uint16) uintptr{
	shmId,_,err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(*sharedMemId), uintptr(size), IPC_CREATE | IPC_EXCL | 0666)
	if err != 0 {
		return uintptr(0)
	}
	SharedMemIdsToShmids[*sharedMemId] = shmId

	accessAddress,_,err := syscall.Syscall(syscall.SYS_SHMAT, shmId, 0, 0)
	if err != 0 {
		return uintptr(0)
	}
	return accessAddress
}

func AccessSharedMemory(sharedMemId *uint16) uintptr {
	shmId,_,err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(*sharedMemId), 0, 0666)
	if err != 0 {
		return uintptr(0)
	}

	accessAddress,_,err := syscall.Syscall(syscall.SYS_SHMAT, shmId, 0, 0)
	if err != 0 {
		return uintptr(0)
	}
	return accessAddress
}

//shm cannot be use after use this method
func CloseShardMemory(accessAddress uintptr) bool {
	_,_,err := syscall.Syscall(syscall.SYS_SHMDT, accessAddress, 0,0)
	if err != 0 {
		return false
	}
	return true
}

func destroySharedMemory(sharedMemId *uint16, force bool) {
	value,ok := SharedMemIdsToShmids[*sharedMemId]
	if ok {
		syscall.Syscall(syscall.SYS_SHMCTL, value, IPC_RMID, 0)
		delete(SharedMemIdsToShmids, *sharedMemId)
	}else if force {
		shmId,_,err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(*sharedMemId), 0, 0666)
		if err != 0 {
			return
		}
		syscall.Syscall(syscall.SYS_SHMCTL, shmId, IPC_RMID, 0)
	}
}