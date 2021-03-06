// Copyright 2014 Jens Rantil. All rights reserved.  Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
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

func Test_Read_ReturnsCharacters_AfterCheckingBOM(t *testing.T) {
	s := "ï»¿Οὐχὶ ταὐτὰ, παρίσταταί μοι, γιγνώσκειν ὦ, ἄνδρες ᾿Αθηναῖοι\n" +
		"ὅταν τ᾿, εἰς τὰ πράγματα ἀποβλέψω, καὶ ὅταν, πρὸς τοὺς\n"

	r := strings.NewReader(s)

	reader := csv.NewReader(r)

	bytes, _ := reader.Read()

	bom := "ï»¿"

	assert.NotEqual(t, bom, bytes[0:3], "first three charactes can never be bom")

}
