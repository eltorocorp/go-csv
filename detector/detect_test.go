package detector

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestIsPotentialDelimiter(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{
			byte('a'),
			false,
		},
		{
			byte('A'),
			false,
		},
		{
			byte('1'),
			false,
		},
		{
			byte('|'),
			true,
		},
		{
			byte('$'),
			true,
		},
	}

	detector := &detector{
		nonDelimiterRegex: regexp.MustCompile(nonDelimiterRegexString),
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, !detector.nonDelimiterRegex.MatchString(string(test.input)))
	}
}

func TestDetectDelimiter1(t *testing.T) {
	detector := New()

	file1, err := os.OpenFile("./Fixtures/test1.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file1.Close()

	file2, err := os.OpenFile("./Fixtures/test2.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file2.Close()

	testCases := []struct {
		r         io.Reader
		delimiter string
	}{
		{file1, ","},
		{file2, ","},
		{strings.NewReader(""), ""},
	}

	for _, tc := range testCases {
		delimiters := detector.DetectDelimiter(tc.r, '"')

		fmt.Println(delimiters)
		assert.Equal(t, []string{tc.delimiter}, delimiters)
	}
}

func TestDetectRowTerminator(t *testing.T) {
	detector := New()

	file, err := os.OpenFile("./Fixtures/test1.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file.Close()

	booString := "boo\r\nhere\r\n\beghosts!\r\n"
	booR := strings.NewReader(booString)

	wooString := "woo\rhere\rbe no pippy!\r"
	wooR := strings.NewReader(wooString)

	emptyR := strings.NewReader("")

	testCases := []struct {
		r          io.Reader
		terminator string
	}{
		{file, "\n"},
		{wooR, "\r"},
		{booR, "\r\n"},
		{badRead{}, ""},
		{emptyR, ""},
	}

	for _, tc := range testCases {
		terminator := detector.DetectRowTerminator(tc.r)
		assert.Equal(t, tc.terminator, terminator)
	}

}

type badRead struct {
}

func (b badRead) Read(p []byte) (int, error) {
	return 0, errors.New("woowoowoo")
}
