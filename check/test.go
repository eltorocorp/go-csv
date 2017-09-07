package main 

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"golang.org/x/text/encoding/french"
	"golang.org/x/text/transform"
)

func main() {
	//string to be transformed 
	s := "voiture"
	fmt.Println(s)


	// Encoding: convert s from UTF-8 to ShiftJIS 
	// declare a bytes.BUffer b and an encoder which will write into the buffer 
	
	var b bytes.Buffer 
	wINUTF8 := transform.NewWriter(&b, french.ShiftJIS.NewEncoder())

	//encode the string 

	wINUTF8.Write([]byte(s))
	wINUTF8.Close()
	fmt.Println(encS)

	// Decoding: convert encs from Shiftjis to UTF8 
	// declare a decoder which reads from the string we have just encoded 

	rInUTF8 := transform.NewReader(strings.NewReader(encS), french.ShiftJIS.NewEncoder())
	
	// decode the string 
	
}
