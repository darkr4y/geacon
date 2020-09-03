package sysinfo

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	Kernel32 = syscall.NewLazyDLL("Kernel32.dll")
	//Advapi32 = syscall.NewLazyDLL("Advapi32.dll")
)

type SystemInfo struct {
	ProcessorArchitecture     ProcessorArchitecture
	Reserved                  uint16
	PageSize                  uint32
	MinimumApplicationAddress uintptr
	MaximumApplicationAddress uintptr
	ActiveProcessorMask       uint64
	NumberOfProcessors        uint32
	ProcessorType             ProcessorType
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

type ProcessorArchitecture uint16

const (
	ProcessorArchitectureAMD64   ProcessorArchitecture = 9
	ProcessorArchitectureARM     ProcessorArchitecture = 5
	ProcessorArchitectureARM64   ProcessorArchitecture = 12
	ProcessorArchitectureIA64    ProcessorArchitecture = 6
	ProcessorArchitectureIntel   ProcessorArchitecture = 0
	ProcessorArchitectureUnknown ProcessorArchitecture = 0xFFFF
)

type ProcessorType uint32

const (
	ProcessorTypeIntel386     ProcessorType = 386
	ProcessorTypeIntel486     ProcessorType = 486
	ProcessorTypeIntelPentium ProcessorType = 586
	ProcessorTypeIntelIA64    ProcessorType = 2200
	ProcessorTypeAMDX8664     ProcessorType = 8664
)

func GetOSVersion() string {
	version, err := syscall.GetVersion()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%d.%d (%d)\n", byte(version), uint8(version>>8), version>>16)

	return fmt.Sprintf("%d.%d.%d\n", byte(version), uint8(version>>8), version>>16)
}

func IsHighPriv() bool {
	token, err := syscall.OpenCurrentProcessToken()
	defer token.Close()
	if err != nil {
		fmt.Printf("open current process token failed: %v\n", err)
		return false
	}
	/*
		ref:
		C version https://vimalshekar.github.io/codesamples/Checking-If-Admin
		Go package https://github.com/golang/sys/blob/master/windows/security_windows.go ---> IsElevated
		maybe future will use ---> golang/x/sys/windows
	*/
	var isElevated uint32
	var outLen uint32
	err = syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&isElevated)), uint32(unsafe.Sizeof(isElevated)), &outLen)
	if err != nil {
		return false
	}
	return outLen == uint32(unsafe.Sizeof(isElevated)) && isElevated != 0
}

func IsOSX64() bool {
	var systemInfo SystemInfo
	fnGetNativeSystemInfo := Kernel32.NewProc("GetNativeSystemInfo")
	if fnGetNativeSystemInfo.Find() != nil {
		panic("not found GetNativeSystemInfo")
	}
	fnGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))
	if (systemInfo.ProcessorArchitecture == ProcessorArchitectureAMD64 ||
		systemInfo.ProcessorArchitecture != ProcessorArchitectureIA64) {
		//x64
		//fmt.Println("amd64")
		return true
	} else
	{
		//x86
		//fmt.Println("386")
		return false
	}
}

func IsProcessX64() bool {
	fnIsWow64Process := Kernel32.NewProc("IsWow64Process")
	//fnIsWow64Process := kernel32.FindProc("IsWow64Process")
	if fnIsWow64Process.Find() != nil {
		panic("not found IsWow64Process")
	}

	is64 := 0

	hProcess, err := syscall.GetCurrentProcess()
	if err != nil {
		panic(err)
	}

	r1, _, err := fnIsWow64Process.Call(uintptr(hProcess), uintptr(unsafe.Pointer(&is64)))
	if r1 == 0 {
		panic(err)
	}
	if is64 == 1 {
		//fmt.Println("procss is x86 (value = 0)")
		return false
	} else {
		//fmt.Println("procss is x64 (value = 1)")
		return true
	}
}

func GetUsername() string {
	username := make([]uint16, 128)
	usernameLen := uint32(len(username)) - 1
	err := syscall.GetUserNameEx(syscall.NameSamCompatible, &username[0], &usernameLen)
	if err != nil {
		panic(err)
	}
	s := syscall.UTF16ToString(username)
	return s
}

func GetCodePageANSI() []byte {
	fnGetACP := Kernel32.NewProc("GetACP")
	if fnGetACP.Find() != nil {
		panic("not found GetACP")
	}
	acp, _, _ := fnGetACP.Call()
	//fmt.Printf("%v\n",acp)
	acpbytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(acpbytes, uint32(acp))
	return acpbytes[:2]

}

func GetCodePageOEM() []byte {
	fnGetOEMCP := Kernel32.NewProc("GetOEMCP")
	if fnGetOEMCP.Find() != nil {
		panic("not found GetOEMCP")
	}
	acp, _, _ := fnGetOEMCP.Call()
	//fmt.Printf("%v\n",acp)
	acpbytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(acpbytes, uint32(acp))
	return acpbytes[:2]
}
