package main 

import (
	"log"
	
	"io/ioutil"
	
	"unicode/utf8"
)

// read the file 

buf, err := ioutil.ReadAll("test3.csv")
if error != nil {
	log.Fatal(err)
}

size := 0
for start := 0; start < len(buf); start += size {
	var r rune 
	if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
		log.Fatalf("invalid utf8 encoding at ofs %d", start)
	}
}




