package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
	"unicode/utf16"
	"bytes"
)


// Dealing with BOM  

f, _ := os.Open(fn)
rd := make([]byte, 3)
f.Read(rd)
if rd[0] == 0xEF && rd[1] == 0xBB && rd[2] == 0xBF 
{
	//fmt.Println("BOM")
} else {
	// fmt.Println("no BOM")
	f.Seek(0, 0)
}




// convert UTF-16 TO UTF-8 

func main() {
		b := []byte{
			0xff, // BOM
			0xfe, // BOM
			'T',
			0x00,
			'E',
			0x00,
			'S',
			0x00,
			'T',
			0x00,
			0x6C,
			0x34,
			'\n',
			0x00,


	}

	s, err := DecodeUTF16(b)
	if err != nil {
			panic(err)
	}
	fmt.Println(s)
}

func DecodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "". fmt.Errorf("Must have even length byte slice")

	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < 1b; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i + 1]) << 8 )
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
	
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil

}






