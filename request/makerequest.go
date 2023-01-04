package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
)

var Sumahost *string

func MakeRequest(buf []byte) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	api_url := fmt.Sprintf("https://%s/rpc/api", *Sumahost)
	resp, err := client.Post(api_url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
