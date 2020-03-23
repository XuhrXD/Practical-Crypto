package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
	"crypto/aes"
)
//Function for Xoring byte arrays.
func XorBlocks(byteArray1 []byte, byteArray2 []byte) []byte {
	xor_result := make([]byte, len(byteArray1))
	for i:=0; i<len(byteArray1); i++ {
		xor_result[i] = byteArray1[i] ^ byteArray2[i]
	}
	return xor_result
}

func encrypt (message []byte, cipherKey []byte, hmacKey []byte, outputfile string) {
	//fmt.Println("ENC Called")
	cipherBlock,err := aes.NewCipher(cipherKey)
	if err != nil {
		fmt.Println("Encryption : Oops! Looks like we have a problem with AES itself.")
		return
	}
	tag := hmac (hmacKey, message)
	//fmt.Println(tag)
	// Appending the orginal message and the HMAC tag
	for i:=0;i<len(tag);i++{
		message = append(message,tag[i])
	}
	// Generating M" = M'| PS
	block_count := len(message)/16
	//fmt.Println("ENC :Block Count", block_count)
	n := len(message)%16
	if n!=0 {
		remaining := 16-n
		padding := byte(remaining)
		for i:=0;i<remaining;i++{
			message = append(message,padding)
		}
	}else if n==0{
		for i:=0;i<16;i++{
			message = append (message,0x10)
		}
	}
	//Generating the cipher text
	iv := make([]byte,16)
	_,err23 := rand.Read(iv)
	if err23 != nil{
		fmt.Println("Cannot Generate an IV")
		os.Exit(1)
	}
	//fmt.Println("Len of IV",iv_len)
	cipherText := make([]byte,16*(block_count+1))//Here, 16 is the AES block size
	xor_result := XorBlocks(iv,message[0:16])
	cipherBlock.Encrypt(cipherText[0:16],xor_result)
	for i:=1;i<=block_count;i++{
		xor_result = XorBlocks(cipherText[((i-1)*16):(i*16)], message[(i*16):((i+1)*16)])
		cipherBlock.Encrypt(cipherText[(16*i):(16*(i+1))],xor_result)
	}
	iv_cipherText := make([]byte,len(iv)+len(cipherText))
	iv_cipherText =  iv
	for j:=0;j<len(cipherText);j++{
		iv_cipherText = append(iv_cipherText,cipherText[j])
	}
	err2 := ioutil.WriteFile(outputfile, iv_cipherText, 0666)
	if err2 !=nil{
		fmt.Println("Outputfile issue")
	}
	fmt.Println(iv_cipherText)
	return
}

func hmac (key []byte,message []byte) [32]byte {
	//fmt.Println("Mac Called")
	blockSize := 64
	ipad := make([]byte,64)
	opad := make([]byte,64)
	//Declaring ipad and opad
	for i:=0;i<64;i++ {
		ipad[i] = 0x36
		opad[i] = 0x5c
	}
	final_key := make([]byte,64)
	// Stuff as per the RFC
	if len(key)>64 {
		temp := sha256.Sum256(key)
		copy(final_key,temp[:])
	}else if len(key)<blockSize {
		final_key = key
		diff := 64 - len(key)
		for i:=0;i<diff;i++{
			final_key = append(final_key,0x00)
		}
	}else if len(key) == 64 {
		final_key = key
	}
	fmt.Println(final_key)
	ipad_final := make([]byte,blockSize)
	opad_final := make([]byte,blockSize)
	//XOR the ipad and opad with final key
	for i:=0;i<blockSize;i++{
		opad_final[i] = opad[i]^final_key[i]
		ipad_final[i] = ipad[i]^final_key[i]
	}
	buff := make([]byte,len(ipad_final)+len(message))
	// ipad and message concatination
	buff = ipad_final
	for i:=0;i<len(message);i++{
		buff = append(buff,message[i])
	}
	hash := sha256.Sum256(buff)
	fmt.Println("Hash_ipad",hash)
	buff2 := make([]byte, len(opad_final)+len(hash))
	buff2 = opad_final
	for i:=0;i<len(hash);i++{
		buff2 = append(buff2,hash[i])
	}
	hash_final := sha256.Sum256(buff2)
	//fmt.Println("Mac Return")
	return hash_final
}

func decrypt (ciphertext []byte, iv []byte, cipherKey []byte, macKey []byte, outputfilename string) {
	//fmt.Println("Enter Decryption mode")
	//fmt.Println(ciphertext)
	cipherBlock,err := aes.NewCipher(cipherKey)
	if err != nil {
		fmt.Println("Decryption : Oops! Looks like we have a problem with AES itself.")
		return
	}
	if (len(ciphertext)%16 != 0){
		fmt.Println("Invalid CipherText")
		os.Exit(1)
	}
	blocks := len(ciphertext)/16
	//fmt.Println("DNC :Block Count", blocks)
	textwithMAC := make([]byte,len(ciphertext))
	cipherBlock.Decrypt(textwithMAC[0:16],ciphertext[0:16])
	temp := XorBlocks(iv,textwithMAC[0:16])
	copy(textwithMAC[0:16],temp[:])
	for i:=1;i<blocks;i++{
		cipherBlock.Decrypt(textwithMAC[16*i:16*(i+1)],ciphertext[16*i:16*(i+1)])
		temp = XorBlocks(ciphertext[16*(i-1) : 16*i], textwithMAC[16*i : 16*(i+1)])
		copy(textwithMAC[16*i : 16*(i+1)], temp[:])
	}
	//fmt.Println(textwithMAC)
	n := textwithMAC[len(textwithMAC)-1]
	for i:= (len(textwithMAC)-1);i>=(len(textwithMAC)-int(n));i--{
		if textwithMAC[i]!=n{
			fmt.Println("INVALID PADDING")
			os.Exit(1)
		}
	}
	//fmt.Println(string(plaintext))
	message_with_tag := make([]byte,len(textwithMAC)-int(n))
	message_with_tag = textwithMAC[:len(textwithMAC)-int(n)]
	plaintext := make([]byte,len(message_with_tag)-32)
	plaintext = textwithMAC[:len(textwithMAC)-int(n)-32]
	tag := make([]byte,32)
	tag = message_with_tag[len(message_with_tag)-32:]
	fmt.Println("TAG ",tag)
	verify := hmac(macKey,plaintext)
	fmt.Println("varify",verify)
	for i:=0;i<32;i++{
		if verify[i]!=tag[i]{
			fmt.Println("INVALID MAC")
			os.Exit(1)
		}
	}
	//Plaintext Output after successfull decryption & MAC verification
	err2 := ioutil.WriteFile(outputfilename, plaintext, 0666)
	if err2 !=nil {
		fmt.Println("Outputfile issue")
	}
	fmt.Println(string(plaintext))
	return
}


func main(){
	if len(os.Args) < 7 {
		fmt.Println("Incorrect Usage")
		os.Exit(1)
	}
	if len(os.Args[3]) !=64 {
		fmt.Println("key length error")
		os.Exit(1)
	}
	key := os.Args[3]
	cipherKey_hex := key[0:32]
	macKey_hex := key[32:64]
	macKey,_ := hex.DecodeString(macKey_hex)
	cipherKey,_ := hex.DecodeString(cipherKey_hex)
	if os.Args[1] == "encrypt" {
		message,_ := ioutil.ReadFile(os.Args[5])
		outputfile := os.Args[7]
		encrypt (message, cipherKey, macKey, outputfile)
		return
	}else if os.Args[1] == "decrypt" {
		rawdata,_ := ioutil.ReadFile(os.Args[5])

		iv := make([]byte,16)
		iv = rawdata[0:16]
		ciphertext := make([]byte,len(rawdata)-16)
		ciphertext = rawdata[16:len(rawdata)]
		outputfile := os.Args[7]
		//fmt.Println(ciphertext)
		decrypt(ciphertext,iv,cipherKey,macKey,outputfile)
		return
	}
}