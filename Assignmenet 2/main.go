package main

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func xor(block1 []byte,block2 []byte) [] byte{
	if len(block1) != len(block2){
		return nil
	}
	result := make([]byte, len(block1))
	for i:=0; i<len(block1); i++ {
		result[i] = block1[i] ^ block2[i]
	}
	return result


}

func encrypt(plaintext []byte,key_enc []byte, key_mac []byte, outputfile string) {

	c, err := aes.NewCipher(key_enc)
	if err != nil {
		fmt.Println(err)
	}
	block_num := len(plaintext)/16
	remain := len(plaintext)%16
	//padding
	block_num++
	remain = 16-remain
	padding := byte(remain)
	for i := 0; i < remain; i++ {
		plaintext = append(plaintext, padding)
	}
	//Generate IV
	iv := make([]byte,16)
	_,iv_err := rand.Read(iv)
	if iv_err == nil {
		fmt.Println("Error when creating IV")
		os.Exit(1)
	}

	xored_cipher := make([]byte,16)
	block_cipher := make([]byte,16)
	ciphertext := make([]byte,16)
	//encrypt first block
	xored_cipher  = xor(plaintext[0:16],iv)
	c.Encrypt(block_cipher,xored_cipher)
	for i:=0; i<len(block_cipher); i++{
		ciphertext[i] = block_cipher[i]
	}

	for i:=0; i < block_num; i++{
		for j:=0; j<len(block_cipher); j++{
			iv[j] = block_cipher[j]
		}
		xored_cipher  = xor(plaintext[16*i:16*(i+1)],iv)
		c.Encrypt(block_cipher,xored_cipher)
		for k:=0; k<len(block_cipher); k++{
			ciphertext = append(ciphertext, block_cipher[k])
		}

	}

	err_write := ioutil.WriteFile(outputfile, cipherText, 0777)
	if err_write !=nil{
		fmt.Println("Can not write to file!")
	}
	return




}

func decrypt(plaintext []byte,iv []byte,key_enc []byte, key_mac []byte, outputfile string) {

}

func main() {
	if len(os.Args) != 9 {
		fmt.Println("Illegal input argument")
		fmt.Println("Expected input is: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -o <outputfile>")
		os.Exit(1)
	} else {
		fmt.Println("good input")
	}
	if len(os.Args[4]) != 64 {
		fmt.Println("key length error")
		os.Exit(1)
	}
	key := os.Args[4]
	enc_str := key[0:32]
	mac_str := key[32:64]
	key_enc, _ := hex.DecodeString(enc_str)
	key_mac, _ := hex.DecodeString(mac_str)
	if os.Args[2] == "encrypt" {
		plaintext, _ := ioutil.ReadFile(os.Args[6])
		outputfile := os.Args[8]
		encrypt(plaintext, key_enc, key_mac, outputfile)
		return
	} else if os.Args[2] == "decrypt" {
		rawdata, _ := ioutil.ReadFile(os.Args[6])
		iv := make([]byte, 16)
		iv = rawdata[0:16]
		ciphertext := make([]byte, len(rawdata)-16)
		ciphertext = rawdata[16:]
		outputfile := os.Args[8]
		fmt.Println(ciphertext)
		decrypt(ciphertext, iv, key_enc, key_mac, outputfile)
		return
	}
}