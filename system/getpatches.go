package system

import (
	"log"
	"sumaprom/auth"
	"sumaprom/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

func (l *ListActiveSystem) Getpatches(Sessionkey *auth.SumaSessionKey) error {
	method := "system.getRelevantErrata"
	type InputParams struct {
		Sessionkey string
		Sid        int
	}

	for i, k := range l.Result {
		inputsparams := InputParams{Sessionkey.Sessionkey, k.Id}

		buf, err := xml.EncodeClientRequest(method, &inputsparams)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("GetUpgPkgs API error: %s\n", err)
		}

		host_patchlist := new(ListAvailablePatches)
		/* bodyB, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Raw logout xml: %s\n", bodyB) */

		err = xml.DecodeClientResponse(resp.Body, host_patchlist)
		if err != nil {
			log.Printf("Decode Getpatches response body failed: %s\n", err)
		}

		l.Result[i].Patches = *host_patchlist

	}

	return nil
}
