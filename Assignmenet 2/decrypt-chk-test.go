package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"


	//"flag"
	"fmt"
	"io/ioutil"
	"os"
	//"io/ioutil"
	//"os"
)
var singleKeySize int = 16
// Encrypt turns plaintext into AES-encrypted ciphertext with CTR mode.  Key
// can be 16, 32 or 64 bytes long.

func showUsage(){
	fmt.Println("***** AES-CTR Encryption/Decryption Tool *******")
	fmt.Println("<mode> -i <input file> -o <outputfile>")
}
func Encrypt(plaintext []byte,key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil{

	}
	//println(err)
	// if err != nil {
	// 	println(err)
	// 	return nil, err
	// }
	println("iv")
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	rand.Read(iv)
	for i:=0;i<len(iv);i++{
		print(iv[i])
		print(" ")
	}
	println(" ")
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	println(ciphertext)
	return ciphertext, nil
}

// Decrypt turns ciphertext using Encrypt into plaintext.
func Decrypt(ciphertext []byte,key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]

	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext, nil
}

func XOrBytes(block []byte,n byte) []byte{
	blockret := make([]byte,len(block))
	for i:=0;i<len(blockret);i++{
		blockret[i] = block[i]^byte(n)
	}
	return blockret
}

func calculateMac(plaintext []byte) []byte{
	sum :=make([]byte,1)

	for i:=0;i<len(plaintext);i++{
		sum[0] = sum[0]+plaintext[i]
	}

	sum[0] = byte(int(sum[0])%256)
	hashedplaintext := append(sum,plaintext...)

	return hashedplaintext
}


func check(ciphertext []byte) string{
	keystring := "0123456789abcdef"
	key := []byte(keystring)
	cipher,err := Decrypt(ciphertext,key)
	if err!=nil{

	}
	res := checkMac(cipher)
	return res
}








func checkMac(textwithmac []byte) string{
	if (len(textwithmac) == 1) && (textwithmac[0] == 0) {
		return "SUCCESS"
	}
	sum :=make([]byte,1)
	for i:=1;i<len(textwithmac);i++{
		sum[0] = sum[0]+textwithmac[i]
	}
	sum[0] = byte(int(sum[0])%256)
	if(sum[0] == textwithmac[0]){
		return "SUCCESS"
	}
	return "INVALID CHECKSUM"
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	return string(content)
}

func writeFile(content []byte, path string) {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	//TODO readfile
	if len(os.Args) < 6 {
		fmt.Println("wrong usage")
		os.Exit(1)

	}
	mode :=  os.Args[1]
	keystring := "0123456789abcdef"
	key := []byte(keystring)
	inputfileName := os.Args[3]
	outputfileName := os.Args[5]
	strbyte, err := ioutil.ReadFile(inputfileName)
	// //TODO simple MAC
	if(err!=nil){
		return
	}
	if mode == "encrypt" {
		hashedplaintext :=calculateMac(strbyte)
		//println(hashedplaintext)
		ciphertext,err:=Encrypt(hashedplaintext,key)
		if(err!=nil){
			return
		}
		writeFile(ciphertext, outputfileName)
		fmt.Println("Encryption Successful!")
	}
	if mode == "decrypt" {
		plaintextwithmac,err := Decrypt(strbyte,key)
		if err != nil{

		}
		fmt.Println("decrypt")
		if checkMac(plaintextwithmac) == "SUCCESS"{
			fmt.Print("success")
			decryptedplaintext :=plaintextwithmac[1:]
			writeFile(decryptedplaintext,outputfileName)
		}


	}


}