package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"io/ioutil"
)

func squareandmultiply (x *big.Int,y *big.Int,n *big.Int,result *big.Int){
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

func prime_factor(a *big.Int) map[string]*big.Int {
	m := make(map[string]*big.Int)

	j:= big.NewInt(2)
	for a.Cmp(big.NewInt(1))==1{
		if big.NewInt(0).Mod(a, j).Cmp(big.NewInt(0)) == 0{
			m[j.String()] = j
			for big.NewInt(0).Mod(a, j).Cmp(big.NewInt(0)) == 0{
				a = big.NewInt(0).Div(a, j)
			}
		}else{
			j = big.NewInt(0).Add(j, big.NewInt(1))
		}
	}
	return m
}

func generatorTest(g,p,order *big.Int,factors map[string]*big.Int)bool{
	var buffer big.Int
	for _,val := range factors{
		buffer1 :=big.NewInt(0).Div(order,val)
		squareandmultiply(g,buffer1,p,&buffer)
		if buffer.Cmp(big.NewInt(1))==0{
			return false
		}
	}
	return true
}

func generateG(p,q *big.Int) (*big.Int){
	order := big.NewInt(0).Sub(p,big.NewInt(1))
	r:= big.NewInt(0).Div(order,q)
	factors := prime_factor(r)
	factors[q.String()] = q
	i := big.NewInt(2)
	for order.Cmp(i)>0{
		if generatorTest(i,p,order,factors){
			return i
		}
		i = new(big.Int).Add(i, big.NewInt(1))
	}
	return order
}

func generateQ()(big.Int){
	var upper_limit big.Int
	exponentiation(big.NewInt(2),big.NewInt(1016),&upper_limit)
	q,_ := rand.Int(rand.Reader,&upper_limit)
	for !q.ProbablyPrime(32){
		buffer_q,err := rand.Int(rand.Reader,&upper_limit)
		q.Set(buffer_q)
		if err!=nil{
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return *q
}

func generateP()(big.Int,big.Int){
	var buffer, upper_limit, lower_limit big.Int
	i := int64(2)
	j := 0
	exponentiation(big.NewInt(2),big.NewInt(1024),&upper_limit)
	exponentiation(big.NewInt(2),big.NewInt(1022),&lower_limit)
	big_1 := big.NewInt(1)
	q := generateQ()
	buffer.Mul(&q,big.NewInt(i))
	p := big.NewInt(0).Add(&buffer,big_1)
	for {
		i++
		j++
		buffer.Mul(&q,big.NewInt(i))
		p_buff := big.NewInt(0).Add(&buffer,big_1)
		p = p_buff
		if ((p.ProbablyPrime(32))&&(p.Cmp(&upper_limit) == -1)&&(p.Cmp(&lower_limit)==1)){
			break
		}
		if j==512{
			q=generateQ()
			i = int64(2)
			j=1
		}
	}
	return *p,q
}

func writesk(p,g,a big.Int,filename string){
	start := "("
	end := ")"
	comma := ","
	p_string := p.String()
	g_string := g.String()
	a_string := a.String()
	final_string :=  start+p_string+comma+g_string+comma+a_string+end
	final_byte := []byte(final_string)
	err := ioutil.WriteFile(filename,final_byte,0644)
	if err !=nil{
		panic(err)
	}
}

func writepk(p,g, g_a big.Int, filename string){
	start := "("
	end := ")"
	comma := ","
	p_string := p.String()
	g_string := g.String()
	g_a_string := g_a.String()
	final_string :=  start+p_string+comma+g_string+comma+g_a_string+end
	final_byte := []byte(final_string)
	err := ioutil.WriteFile(filename,final_byte,0644)
	if err !=nil{
		panic(err)
	}

}

func main() {
	if len(os.Args)!=3{
		fmt.Println("Usage Error")
		os.Exit(1)
	}
	filename_pk:=os.Args[1]
	filename_sk := os.Args[2]
	p,q:=generateP()
	g := generateG(&p,&q)
	var g_a big.Int
	upper_limit := big.NewInt(0).Sub(&p,big.NewInt(1))
	a,_ := rand.Int(rand.Reader,upper_limit)
	squareandmultiply(g,a,&p,&g_a)
	writepk(p,*g,g_a,filename_pk)
	writesk(p,*g,*a,filename_sk)
}

