package sysinfo

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

func GetOSVersion() []byte {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, _ := cmd.CombinedOutput()
	return out
}

func IsHighPriv() bool {
	fd, err := os.Open("/root")
	defer fd.Close()
	if err != nil {
		return false
	}
	return false
}

func IsOSX64() int {
	cmd := exec.Command("sysctl", "hw.cpu64bit_capable")
	out, _ := cmd.CombinedOutput()
	out = bytes.ReplaceAll(out, []byte("hw.cpu64bit_capable: "), []byte(""))
	if string(out) == "1" {
		return 1
	}
	return 0
}

func IsProcessX64() int {
	if runtime.GOARCH == "amd64" {
		return 0
	}
	return 1
}

func GetCodePageANSI() []byte {
	//hardcode for test
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 936)
	return b
}

func GetCodePageOEM() []byte {
	//hardcode for test
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 936)
	return b
}

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	usr := user.Username
	if IsHighPriv() {
		usr = usr + " *"
	}
	return usr
}
