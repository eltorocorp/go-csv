// Copyright 2014 Jens Rantil. All rights reserved.  Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package csv

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)

func TestUnReader(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a,b,c\n")
	r := newUnreader(b)
	if ru, _, _ := r.ReadRune(); ru != 'a' {
		t.Error("Unexpected char:", ru, "Expected:", 'a')
	}
	if ok, _ := r.NextIsString(",b,c"); !ok {
		t.Error("Unexpected next string.")
	}
	r.UnreadRune('d')
	if ok, _ := r.NextIsString("d,b,c"); !ok {
		t.Error("Unreading failed.")
	}
	if ok, _ := r.NextIsString("b,c"); ok {
		t.Error("Unexpected next string.")
	}
	if ru, _, _ := r.ReadRune(); ru != 'd' {
		t.Error("Unexpected char:", ru, "Expected:", 'd')
	}
}

func testReadingSingleLine(t *testing.T, r *Reader, expected []string) error {
	record, err := r.Read()
	if c := len(record); c != len(expected) {
		t.Fatal("Wrong number of fields:", c, "Expected:", len(expected))
	}
	if !reflect.DeepEqual(record, expected) {
		t.Error("Incorrect records.")
		t.Error(record)
		t.Error(expected)
	}
	return err
}

func TestReadingSingleFieldLine(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a\n")
	r := NewReader(b)

	err := testReadingSingleLine(t, r, []string{"a"})
	if err != nil && err != io.EOF {
		t.Error("Unexpected error:", err)
	}
}

func TestReadingSingleLine(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a b c\n")
	r := NewReader(b)

	err := testReadingSingleLine(t, r, []string{"a", "b", "c"})
	if err != nil && err != io.EOF {
		t.Error("Unexpected error:", err)
	}
}

func TestReadingTwoLines(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a b c\nd e f\n")
	r := NewReader(b)
	err := testReadingSingleLine(t, r, []string{"a", "b", "c"})
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	err = testReadingSingleLine(t, r, []string{"d", "e", "f"})
	if err != nil && err != io.EOF {
		t.Error("Expected EOF, but got:", err)
	}
}

func TestReadingBasicCommaDelimitedFile(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("\"b\"\n")
	r := NewReader(b)

	err := testReadingSingleLine(t, r, []string{"b"})
	if err != nil && err != io.EOF {
		t.Error("Unexpected error:", err)
	}
}

func TestReadingCommaDelimitedFile(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a \"b\" c\n")
	r := NewReader(b)

	err := testReadingSingleLine(t, r, []string{"a", "b", "c"})
	if err != nil && err != io.EOF {
		t.Error("Unexpected error:", err)
	}
}

func TestReadAll(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)
	b.WriteString("a \"b\" c\nd e \"f\"\n")
	r := NewReader(b)

	data, err := r.ReadAll()
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	equals := reflect.DeepEqual(data, [][]string{
		{
			"a",
			"b",
			"c",
		},
		{
			"d",
			"e",
			"f",
		},
	})
	if !equals {
		t.Error("Unexpected output:", data)
	}
}

func testReaderQuick(t *testing.T, quoting int) {
	f := func(records [][]string, doubleQuote bool, escapeChar, del, quoteChar rune, lt string) bool {
		dialect := Dialect{
			Quoting:        quoting,
			EscapeChar:     escapeChar,
			QuoteChar:      quoteChar,
			Delimiter:      del,
			LineTerminator: lt,
		}
		if doubleQuote {
			dialect.DoubleQuote = DoDoubleQuote
		} else {
			dialect.DoubleQuote = NoDoubleQuote
		}
		b := new(bytes.Buffer)
		w := NewDialectWriter(b, dialect)
		w.WriteAll(records)

		r := NewDialectReader(b, dialect)
		data, err := r.ReadAll()
		if err != nil {
			t.Error("Error when reading CSV:", err)
			return false
		}

		equal := reflect.DeepEqual(records, data)
		if !equal {
			t.Error("Not equal:", records, data)
		}
		return equal
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Test writing to and then reading from using various CSV dialects.
func TestReaderQuick(t *testing.T) {
	t.Parallel()

	testWriterQuick(t, QuoteAll)
	testWriterQuick(t, QuoteMinimal)
	testWriterQuick(t, QuoteNonNumeric)
}

// Test fot UTF8 Support

func Test_Read_UTF8_ReturnsF(t *testing.T) {
	text := "Οὐχὶ ταὐτὰ, παρίσταταί μοι, γιγνώσκειν ὦ, ἄνδρες ᾿Αθηναῖοι\n" +
		"ὅταν τ᾿, εἰς τὰ πράγματα ἀποβλέψω, καὶ ὅταν, πρὸς τοὺς\n" +
		"λόγους, οὓςm ἀκούω·, τοὺς μὲν γὰρ, λόγους περὶ τοῦ\n"
	b := NewDialectReader(strings.NewReader(text), Dialect{
		Delimiter:      ',',
		LineTerminator: "\n",
	})
	line, _ := b.Read()
	result := reflect.DeepEqual(line[0], "Οὐχὶ ταὐτὰ")
	if !result {
		t.Error("Unexpected output:", line[0])
	}
}

func Test_Read_UTF8_Properly_ReadsCharacters(t *testing.T) {
	// Test data.
	text := "Οὐχὶ ταὐτὰ, παρίσταταί μοι, γιγνώσκειν ὦ, ἄνδρες ᾿Αθηναῖοι\n" +
		"ὅταν τ᾿, εἰς τὰ πράγματα ἀποβλέψω, καὶ ὅταν, πρὸς τοὺς\n" +
		"λόγους, οὓςm ἀκούω·, τοὺς μὲν γὰρ, λόγους περὶ τοῦ\n"
	// Create reader
	r := NewDialectReader(strings.NewReader(text), Dialect{
		Delimiter:      ',',
		LineTerminator: "\n",
	})
	// Ignore first two lines.
	r.Read()
	r.Read()
	// Read the third line.
	line, _ := r.Read()
	// Check result.
	result := reflect.DeepEqual(line[2], " τοὺς μὲν γὰρ")
	// Verify the result is as expected, if not fail.
	if !result {
		t.Error("Unexpected output:", line[2])
	}
}

func Test_Read_UTF16_ReadsCharacters(t *testing.T) {
	text := "楆獲ⱴ慌瑳䄬摤敲, 獳ⰱ摁牤獥㉳䌬瑩ⱹ瑓瑡ⱥ楚Ɒ楚㑰䄊䅄\n" +
		"ⱍ䕗呓㔬㠱, 䥖䱌䝁⁅噁ⱅ䰬协䄠䝎䱅卅䌬ⱁ〹㄰ⰶ㈵㘰䈊䉏奂䠬䱉ⱌㄳ‱\n"

	r := NewDialectReader(strings.NewReader(text), Dialect{
		Delimiter:      ',',
		LineTerminator: "\n",
	})

	r.Read()
	line, _ := r.Read()

	result := reflect.DeepEqual(line[0], "ⱍ䕗呓㔬㠱")

	if !result {
		t.Error("Unexpected result:", line[1])
	}
}

func Test_Read_UTF8_ReadsCharacters(t *testing.T) {
	//Test data
	text := "∮ E⋅da = Q,  n → ∞, ∑ f(i) = ∏ g(i)\n" +
		"∀x∈ℝ: ⌈x⌉ = −⌊−x⌋, α ∧ ¬β = ¬(¬α ∨ β)\n"

	r := NewDialectReader(strings.NewReader(text), Dialect{
		Delimiter:      ',',
		LineTerminator: "\n",
	})

	r.Read()

	line, _ := r.Read()

	result := reflect.DeepEqual(line[1], " α ∧ ¬β = ¬(¬α ∨ β)")

	if !result {
		t.Error("Unexpected output:", line[1])
	}

}

func Test_UnicodeBOM_ReadCharacters(t *testing.T) {
	// So the characters are indexed at 0,2 and 4 because Unicode characters can take more than one position
	text := "ï»¿"

	bomFound := true
	fmt.Println(len(text))
	for index, value := range text {
		fmt.Println(index)
		switch index {
		case 0:
			if byte(value) == 0xEF {
				bomFound = bomFound && true
			} else {
				bomFound = bomFound && false
			}

		case 2:
			if byte(value) == 0xBB {

				bomFound = bomFound && true
			} else {

				bomFound = bomFound && false

			}

		case 4:
			if byte(value) == 0xBF {
				bomFound = bomFound && true
			} else {
				bomFound = bomFound && false
			}
		}
		fmt.Println(value, byte(value), bomFound)
		fmt.Printf("%X\n", byte(value))

	}

	if !bomFound {
		t.Error("Unexpected output:", bomFound)
	}
}

func Test_UnicodeBOMUTF16BE_ReadCharacters(t *testing.T) {

	text := "þÿ"

	bomFound := true

	for index, value := range text {
		switch index {
		case 0:
			if byte(value) == 0xFE {
				bomFound = bomFound && true
			} else {
				bomFound = bomFound && false
			}

		case 1:
			if byte(value) == 0xFF {
				bomFound = bomFound && true
			} else {
				bomFound = bomFound && false

			}

		}

		fmt.Println(value, byte(value), bomFound)
		fmt.Printf("%X\n", byte(value))
	}

	if !bomFound {
		t.Error("Unexpected output:", bomFound)
	}
}

func Test_Read_UnicodeBOM4_ReadCharacters(t *testing.T) {
	//
	text := "Οὐχὶ ταὐτὰ, παρίσταταί μοι, γιγνώσκειν ὦ, ἄνδρες ᾿Αθηναῖοι\n" +
		"ὅταν τ᾿, εἰς τὰ πράγματα ἀποβλέψω, καὶ ὅταν, πρὸς τοὺς\n"

	r := NewDialectReader(strings.NewReader(text), Dialect{
		Delimiter:      ',',
		LineTerminator: "\n",
	})

	r.Read()

	line, _ := r.Read()

	result := reflect.DeepEqual(line[1], "Οὐχὶ ταὐτὰ")

	if !result {
		t.Error("Unexpected output:", line[1])
	}

}

// \xEF\xBB\

func Test_Read_Unicode(t *testing.T) {
	s, err := "Οὐχὶ ταὐτὰ, παρίσταταί μοι, γιγνώσκειν ὦ, ἄνδρες ᾿Αθηναῖοι\n" +
		"ὅταν τ᾿, εἰς τὰ πράγματα ἀποβλέψω, καὶ ὅταν, πρὸς τοὺς\n"

	bom := [3]byte

	_, err = io.ReadFull(s, bom[:])
	if err != nil {
		log.Fatal(err)
	}

	if bom[0] != 0xef || bom[1] != 0xbb || bom[2] != 0xbf {
		_, err = s.seek(0, 0) //not bom seek back to beginning
		if err != nil {
			log.Fatal(err)
		}
	}

	//runes := []rune(s)

	//buf := []byte(s)

	//bomFound := true

	//fmt.Println(utf8.RuneStart(buf[1]))
	//fmt.Println(utf8.RuneStart(buf[2]))
	//fmt.Println(utf8.RuneStart(buf[3]))
	//fmt.Println(utf8.RuneStart(buf[4]))
	//fmt.Println(utf8.RuneStart(buf[5]))

	//fmt.Printf("%c", runes[0])
	//fmt.Printf("%c", runes[1])
	//fmt.Printf("%c", runes[2])
	//fmt.Printf("%c", runes[3])
	//fmt.Printf("%c", runes[4])
	//fmt.Printf("%c", runes[5])
	//fmt.Printf("%c", runes[6])

	//for index, value := range runes {
	//switch index {
	//case 0:
	//if byte(value) == 0xEF {
	//bomFound = bomFound && true
	//} else {
	//bomFound = bomFound && false
	//}

	//case 1:
	//if byte(value) == 0xBB {

	//bomFound = bomFound && true
	//} else {

	//bomFound = bomFound && false

	//}

	//case 2:
	//if byte(value) == 0xBF {
	//bomFound = bomFound && true
	//} else {
	//bomFound = bomFound && false
	//}

	//}

	//fmt.Println(value, byte(value))

	//r := NewDialectReader(strings.NewReader(s), Dialect{
	//Delimiter:      ',',
	//LineTerminator: "\n",
	//})

	//r.readField()

	//line, _ := r.readField()

	//fmt.Printf(line)

	//if !bomFound {
	//t.Error("Unexpected output:", bomFound)
	//	}

	//}

}
