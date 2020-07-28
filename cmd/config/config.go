package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
	AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
-----END PUBLIC KEY-----`)
	RsaPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
	AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA	
-----END PRIVATE KEY-----`)

	C2        = "127.0.0.1:80"
	plainHTTP = "http://"
	sslHTTP   = "https://"
	GetUrl    = plainHTTP + C2 + "/load"
	PostUrl   = plainHTTP + C2 + "/submit.php?id="
	WaitTime  = 10000 * time.Millisecond

	IV        = []byte("abcdefghijklmnop")
	GlobalKey []byte
	AesKey    []byte
	HmacKey   []byte
	Counter   = 0
)

const (
	DebugMode = true
)
