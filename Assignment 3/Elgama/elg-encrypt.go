package main
 import(
 	"crypto/aes"
 	"crypto/cipher"
 	"crypto/rand"
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

func concatinateandHash (g_a,g_b,g_ab big.Int)([]byte){
	g_a_string := g_a.String()
	g_b_string := g_b.String()
	g_ab_string := g_ab.String()
	buffer := g_a_string+" "+g_b_string+" "+g_ab_string
	hash := sha256.Sum256([]byte(buffer))
	return hash[:]
}
func encrypt(key []byte, plaintext string)([]byte) {
	block,err := aes.NewCipher(key)
	if err !=nil{
		fmt.Println("AES encryption error")
		os.Exit(1	)
	}
	aesgcm,_ := cipher.NewGCMWithNonceSize(block,16)
	iv := make([]byte,16)
	_,_ = rand.Read(iv)
	ciphertext := aesgcm.Seal(nil,iv,[]byte(plaintext),nil)
	result := make([]byte,len(ciphertext)+len(iv))
	result = iv
	for j:=0;j<len(ciphertext);j++{
		result = append(result,ciphertext[j])
	}

	return result

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

func main() {
	var g_b,g_ab big.Int
	if len(os.Args)!=4{
		fmt.Println("Usage error")
		os.Exit(1)
	}
	plaintext := os.Args[1]
	pk := os.Args[2]
	ciphertext_file := os.Args[3]
	file_content,err := os.Open(pk)
	if err!=nil{
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	defer file_content.Close()
	file_content_string,_ := ioutil.ReadAll(file_content)
	p,g,g_a := textSeperator([]byte(file_content_string))
	upper_limit := big.NewInt(0).Sub(p,big.NewInt(1))
	b , _ := rand.Int(rand.Reader,upper_limit)
	squareandmultiply(g,b,p,&g_b)
	squareandmultiply(g_a,b,p,&g_ab)
	k := concatinateandHash(*g_a,g_b,g_ab)
	ciphtertext :=  encrypt(k,plaintext)
	ciphertext_hex := make([]byte,hex.EncodedLen(len(ciphtertext)))
	hex.Encode(ciphertext_hex,ciphtertext)
	writetofile(g_b,ciphertext_hex,ciphertext_file)
}