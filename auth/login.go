package auth

import (
	"log"
	"sumaprom/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

type Sumalogin struct {
	Login  string `xmlrpc:"username"`
	Passwd string `xmlrpc:"password"`
}

type SumaSessionKey struct {
	Sessionkey string
}

type SumaLogout struct {
	ReturnInt int
}

func Login(method string, args Sumalogin) (reply SumaSessionKey, err error) {
	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Login Request error: %s\n", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Login SUMA API failed. Status code: %d\n", resp.StatusCode)
	}

	err = xml.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		log.Fatalf("Decode Login response body failed: %s\n", err)
	}
	if resp.StatusCode == 200 && reply.Sessionkey != "" {
		log.Println("Login successful.")
	}
	//fmt.Printf("xml decoded %#v\n", reply)
	return
}
