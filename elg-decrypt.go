package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"encoding/hex"
)

func squareandmultiply (x *big.Int,y *big.Int,n *big.Int,result *big.Int){//result = x^y mod n
	var p,r big.Int
	p.Set(y)
	r.Set(x)
	buffer := big.NewInt(1)
	result.Set(buffer)
	for p.BitLen()>0{
		if p.Bit(0)!=0 {
			result.Mul(result,&r)
			result.Mod(result,n)
		}
		p.Rsh(&p,1)
		r.Mul(&r,&r)
		r.Mod(&r,n)
	}
}

func concatinateandHash (g_a,g_b,g_ab big.Int)([]byte){
	g_a_string := g_a.String()
	g_b_string := g_b.String()
	g_ab_string := g_ab.String()
	buffer := g_a_string+" "+g_b_string+" "+g_ab_string
	hash := sha256.Sum256([]byte(buffer))
	return hash[:]
}

func writetofile(g_b big.Int,input[]byte,filename string){
	start := "("
	end := ")"
	comma := ","
	g_b_string := g_b.String()
	final_string :=  start+g_b_string+comma+string(input)+end
	final_byte := []byte(final_string)
	err := ioutil.WriteFile(filename,final_byte,0644)
	if err !=nil{
		panic(err)
	}
}

func textSeperator(b []byte)(*big.Int,*big.Int,*big.Int){
	var p,g,g_a big.Int
	file := string(b)
	file = strings.Replace(file, " ", "", -1)
	file_string := strings.Split(file,",")
	file_string_p := strings.Split(file_string[0],"(")
	file_string_g :=file_string[1]
	file_string_g_a := strings.Split(file_string[len(file_string)-1],")")
	p.SetString(file_string_p[1],10)
	g.SetString(file_string_g,10)
	g_a.SetString(file_string_g_a[0],10)
	return &p,&g,&g_a
}

func ciphertextSeperator(b []byte)(big.Int,string){
	var g_b big.Int
	file := string(b)
	file = strings.Replace(file, " ", "", -1)
	file_string := strings.Split(file,",")
	file_string_g_b := strings.Split(file_string[0],"(")
	iv_cipherText := file_string[1][:len(file_string[1])-1]
	g_b.SetString(file_string_g_b[1],10)
	return g_b,iv_cipherText
}

func decrypt(key []byte, ciphertextwithiv []byte){
	block,err := aes.NewCipher(key)
	if err !=nil{
		fmt.Println("AES encryption error")
		os.Exit(1)
	}
	aesgcm,_ := cipher.NewGCMWithNonceSize(block,16)
	ciphertext := make([]byte, len(ciphertextwithiv)-16)
	iv := make([]byte,16)
	copy(iv,ciphertextwithiv[0:16])
	copy(ciphertext,ciphertextwithiv[16:])
	plaintext,err := aesgcm.Open(nil,iv,ciphertext,nil)
	if err!=nil{
		fmt.Println("AES_GCM decryption error")
		os.Exit(1)
	}

	fmt.Println(string(plaintext))
}

func main(){
	var g_b,g_a,g_ab big.Int
	if len(os.Args)!=3{
		fmt.Println("Usage error")
		os.Exit(1	)
	}
	ciphertext := os.Args[1]
	secret_key := os.Args[2]
	file_content,err := os.Open(secret_key)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	secret_key_string,_ := ioutil.ReadAll(file_content)
	p,g,a := textSeperator([]byte(secret_key_string))
	file_content,err = os.Open(ciphertext)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	ciphertext_string,_ := ioutil.ReadAll(file_content)
	g_b,iv_cipherText_hex:= ciphertextSeperator(ciphertext_string)
	squareandmultiply(g,a,p,&g_a)
	squareandmultiply(&g_b,a,p,&g_ab)
	key :=concatinateandHash(g_a,g_b,g_ab)
	ciphertextwithiv := make([]byte,hex.DecodedLen(len(iv_cipherText_hex)))
	hex.Decode(ciphertextwithiv,[]byte(iv_cipherText_hex))
	decrypt(key,ciphertextwithiv)
}