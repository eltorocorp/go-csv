package main 

import (
	"bytes"
	"fmt"
	
	"golang.org/x/text/transform"
	
	
	
	
	
	
)

func main() {
 
	// file to be transformed 

	f, _ := ("test3.csv")
	
	fmt.Println(f)

	// encoding file from UTF16 TO UTF8 

	var b bytes.Buffer
	wINUTF16 := transform.NewWriter(&b, UTF8.NewEncoder())

	// encode the file 

	wINUTF16.Write([]byte(f))
	wINUTF16.Close()

	// print encoded bytes 

	fmt.Printf("%#v\n", b)
	encF := b.csv()

	fmt.Println(encF)
	
}

