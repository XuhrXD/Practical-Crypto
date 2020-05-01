package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
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

func exponentiation(num *big.Int,exp *big.Int,result *big.Int){
	var p,r big.Int
	p.Set(exp)
	r.Set(num)
	buffer := big.NewInt(1)
	result.Set(buffer)
	for p.BitLen()>0{
		if p.Bit(0)!=0 {
			result.Mul(result,&r)

		}
		p.Rsh(&p,1)
		r.Mul(&r,&r)
	}
}
func textSeperator_one(b []byte)(*big.Int){
	var g_b big.Int
	text := string(b)
	buffer := strings.Split(text,"(")
	final_text :=strings.Split(buffer[len(buffer)-1],")")
	g_b.SetString(final_text[0],10)
	return &g_b
}

func textSeperator(b []byte)(*big.Int,*big.Int,*big.Int){
	var p,g,g_a big.Int
	file := string(b)
	file_string := strings.Split(file,",")
	file_string_p := strings.Split(file_string[0],"(")
	file_string_g :=file_string[1]
	file_string_g_a := strings.Split(file_string[len(file_string)-1],")")
	p.SetString(file_string_p[1],10)
	g.SetString(file_string_g,10)
	g_a.SetString(file_string_g_a[0],10)
	return &p,&g,&g_a
}

func writetomyfile(g_b *big.Int,filename string){
	start := "("
	end := ")"
	g_b_string := g_b.String()
	final_string :=  start+g_b_string+end
	final_byte := []byte(final_string)
	err := ioutil.WriteFile(filename,final_byte,0644)
	if err !=nil{
		panic(err)
	}
}

func main(){
	var secret big.Int
	if len(os.Args)!=3{
		fmt.Println("Usage Error")
		os.Exit(1)
	}
	filename_from_bob := os.Args[1]
	secret_key := os.Args[2]
	file_content,err := os.Open(filename_from_bob)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	file_content_string,_ := ioutil.ReadAll(file_content)
	g_b := textSeperator_one([]byte(file_content_string))
	secret_key_content,err := os.Open(secret_key)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer secret_key_content.Close()
	secret_string,_ := ioutil.ReadAll(secret_key_content)
	p,_,a := textSeperator([]byte(secret_string))
	squareandmultiply(g_b,a,p,&secret)
	fmt.Println(&secret)
}