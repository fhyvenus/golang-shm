package shm

import (
	"log"
	_ "syscall"
	"unsafe"
)

type nameStr struct {
	name []byte
	age int
}

func MakeShareMemory(shmId uint32) uintptr {
	//string方式
	//shareName,_ := syscall.BytePtrFromString("Global/1")
	//shareNameId := (*uint16)(unsafe.Pointer(shareName))

	//数值方式
	smid := (*uint16)(unsafe.Pointer(&shmId))
	segment := CreateSharedMemory(1024, smid)
	if segment == 0 {
		//under linux,segment may have not been destroyed if the MS crashed:destroy
		//See /proc/sysvipc/shm for the list of allocated segments
		//如果是重新直接使用的话，这边不能做销毁操作，应该直接使用共享内存，快速启动进程
		destroySharedMemory(smid, true)
		//util.InfoF("destroyed shared memory segment,smid:%v", smid)

		segment = CreateSharedMemory(1024, smid)
	}

	if segment == 0 {
		//util.InfoF("cannot create shared memory segment,smid:%v", smid)
		return uintptr(0)
	}
	//defer CloseShardMemory(segment)

	return segment
}

func ReadSharedMemory(shmId uint32) uintptr{
	smid := (*uint16)(unsafe.Pointer(&shmId))
	return AccessSharedMemory(smid)
}

func TestMake() {
	segment := MakeShareMemory(1)
	if segment == 0 {
		return
	}

	nameList := (*nameStr)(unsafe.Pointer(segment))
	nameList.name = []byte("chy1")
	nameList.age = 1111
	log.Printf("make:%v", segment)

	//*nameList = append(*nameList, &nameStr{name:[]byte("chy1"),age:30})
	//*nameList = append(*nameList, &nameStr{name:[]byte("chy2"),age:31})
	//*nameList = append(*nameList, &nameStr{name:[]byte("chy13"),age:32})
}

func TestRead() {
	segment := ReadSharedMemory(1)
	if segment == 0 {
		return
	}

	nameList := (*nameStr)(unsafe.Pointer(segment))
	log.Printf("namelist:%v %v", string(nameList.name), nameList.age)
	/*
		if nameList != nil {
			for key,_ :=range *nameList {
				log.Printf("name:%v age:%v", string((*nameList)[key].name), (*nameList)[key].age)
			}
		}

	*/
}