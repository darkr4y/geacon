package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"geacon/cmd/config"
)

const HmacHashLen = 16

func PaddingWithA(rawData []byte) []byte {
	newBuf := bytes.NewBuffer(rawData)
	step := 16
	for pad := newBuf.Len() % step; pad < step; pad++  {
		newBuf.Write([]byte("A"))
	}
	return newBuf.Bytes()
}

func AesCBCEncrypt(rawData,key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	rawData = PaddingWithA(rawData)
	cipherText := make([]byte,blockSize+len(rawData))
	mode := cipher.NewCBCEncrypter(block,config.IV)
	mode.CryptBlocks(cipherText[blockSize:],rawData)
	return cipherText, nil
}

func AesCBCDecrypt(encryptData, key []byte) ([]byte,error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	if len(encryptData) % blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, config.IV)
	mode.CryptBlocks(encryptData, encryptData)
	return encryptData,nil
}

func HmacHash(encrytedBytes []byte) []byte {
	hmacEntry := hmac.New(sha256.New, config.HmacKey)
	hmacEntry.Write(encrytedBytes)
	return hmacEntry.Sum(nil)[:16]
}