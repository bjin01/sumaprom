package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func (d *Sumaconf) Decrypt_Sumaconf(conf_file *string) Sumaconf {
	var key = []byte("A8CA6E5155A14BB0D8F4E7CE0E23B8A2")
	yfile, err := ioutil.ReadFile(*conf_file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yfile, d)
	if err != nil {
		log.Fatal(err)
	}

	return Sumaconf{
		Server:   d.Server,
		Userid:   d.Userid,
		Password: Decrypt(key, d.Password),
	}

}

func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
