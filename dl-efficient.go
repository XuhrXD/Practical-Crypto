package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"strings"
)

func squareandmultiply (a,b,c big.Int)(big.Int){
	var p,r big.Int
	x := a
	y := b
	n := c
	var result big.Int
	p.Set(&y)
	r.Set(&x)
	buffer := big.NewInt(1)
	result.Set(buffer)
	for p.BitLen()>0{
		if p.Bit(0)!=0 {
			result.Mul(&result,&r)
			result.Mod(&result,&n)
		}
		p.Rsh(&p,1)
		r.Mul(&r,&r)
		r.Mod(&r,&n)
	}

	return result
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

func new_xab(x,a,b *big.Int ,p,g,h big.Int){
	remain := big.NewInt(0).Mod(x,big.NewInt(3))
	var p_c big.Int
	p_c.Set(&p)
	p_c.Sub(&p_c,big.NewInt(1))
	if remain.Cmp(big.NewInt(0))==0 {
		x.Mul(x,x)
		x.Mod(x,&p)
		a.Mul(a,big.NewInt(2))
		a.Mod(a,&p_c)
		b.Mul(b,big.NewInt(2))
		b.Mod(b,&p_c)
	}
	if remain.Cmp(big.NewInt(1))==0 {
		x.Mul(x,&g)
		x.Mod(x,&p)
		a.Add(a,big.NewInt(1))
		a.Mod(a,&p_c)
	}
	if remain.Cmp(big.NewInt(2))==0 {
		x.Mul(x,&h)
		x.Mod(x,&p)
		b.Add(b,big.NewInt(1))
		b.Mod(b,&p_c)
	}
}

func pollard(p,g,h big.Int){
	x := big.NewInt(1)
	a := big.NewInt(0)
	b := big.NewInt(0)
	X := big.NewInt(1)
	A := big.NewInt(0)
	B := big.NewInt(0)
	i:=big.NewInt(1)
	for {
		new_xab(x,a,b,p,g,h)
		new_xab(X,A,B,p,g,h)
		new_xab(X,A,B,p,g,h)
		if x.Cmp(X) == 0{
			var m,r big.Int
			//fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n",i,x,a,b,X,A,B)
			m.Sub(A,a)
			r.Sub(b,B)
			r.Mod(&r,&p)
			if r.Cmp(big.NewInt(0))==0{
				fmt.Println("Faliure")
				return
			}
			var p_copy, reverse , result big.Int
			p_copy.Set(&p)
			p_copy.Sub(&p_copy,big.NewInt(2))
			reverse = squareandmultiply(r,p_copy,p)
			result.Mul(&reverse,&m)
			result.Mod(&result,&p)
			fmt.Println(&result)
			return
		}
		i.Add(i,big.NewInt(1))
	}

}

func baby_giant(p,g,h big.Int){
	var power,y big.Int
	var i,j int32
	buffer := big.NewInt(0).Sub(&p,big.NewInt(1))
	buffer.Sqrt(buffer)
	n := int32(math.Ceil(float64(buffer.Int64())))+1
	m:=make(map[string]int32)
	for i=0;i<n;i++{
		index := squareandmultiply(g,*big.NewInt(int64(i)),p)
		m[index.String()]=i
	}
	p_2:=big.NewInt(0).Sub(&p,big.NewInt(int64(2)))
	power.Mul(big.NewInt(int64(n)),p_2)
	c := squareandmultiply(g,power,p)
	for j=0;j<n;j++{
		big_j := big.NewInt(int64(j))
		buffer:= squareandmultiply(c,*big_j,p)
		buffer1 := big.NewInt(0).Mul(&h,&buffer)
		y.Mod(buffer1,&p)
		if val,ok := m[y.String()];ok{
			r := int32(j*n)
			result := int32(r+val)
			fmt.Println(result)
			return

		}
	}

}

func main(){
	if len(os.Args)!=2{
		fmt.Println("Usage Error")
		os.Exit(1)
	}
	filename := os.Args[1]
	file_content,err := os.Open(filename)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	file_content_string,_ := ioutil.ReadAll(file_content)
	p,g,h := textSeperator([]byte(file_content_string))
	//pollard(*p,*g,*h)
	baby_giant(*p,*g,*h)
}