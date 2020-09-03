package sysinfo

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
)

func arrayToString(x [65]int8) string {
	var buf [65]byte
	for i, b := range x {
		buf[i] = byte(b)
	}
	str := string(buf[:])
	if i := strings.Index(str, "\x00"); i != -1 {
		str = str[:i]
	}
	return str
}

func getUname() syscall.Utsname {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		fmt.Printf("Uname: %v", err)
		return syscall.Utsname{} //nil
	}
	return uname
}

func GetOSVersion() string {
	uname := getUname()

	if len(uname.Release) > 0 {
		return arrayToString(uname.Release)
	}
	return "0.0"
}

func IsHighPriv() bool {
	fd, err := os.Open("/root")
	defer fd.Close()
	if err != nil {
		return false
	}
	return true
}

func IsOSX64() bool {
	uname := getUname()
	if arrayToString(uname.Machine) == "x86_64" {
		return true
	}
	return false
}

func IsProcessX64() bool {
	if runtime.GOARCH == "amd64" {
		return false
	}
	return true
}

func GetCodePageANSI() []byte {
	//hardcode for test
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, 936)
	return b
}

func GetCodePageOEM() []byte {
	//hardcode for test
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, 936)
	return b
}

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	usr := user.Username
	return usr
}
