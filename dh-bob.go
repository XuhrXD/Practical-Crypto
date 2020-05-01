package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"crypto/rand"
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
	var for_alice,secret big.Int
	if len(os.Args)!=3{
		fmt.Println("Usage Error")
		os.Exit(1)
	}
	filename_from_alice := os.Args[1]
	filename_back_to_alice := os.Args[2]
	file_content,err := os.Open(filename_from_alice)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	file_content_string,_ := ioutil.ReadAll(file_content)
	p,g,g_a := textSeperator([]byte(file_content_string))
	upper_limit := big.NewInt(0).Sub(p,big.NewInt(1))
	b,_ := rand.Int(rand.Reader,upper_limit)
	squareandmultiply(g,b,p,&for_alice)
	writetomyfile(&for_alice,filename_back_to_alice)
	squareandmultiply(g_a,b,p,&secret)
	fmt.Println(&secret)
}
