package detector

import (
	"os"
	"regexp"
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

	file, err := os.OpenFile("./Fixtures/test1.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file.Close()

	delimiters := detector.DetectDelimiter(file, '"')
	fmt.Println(delimiters)
	assert.Equal(t, []string{","}, delimiters)
}

func TestDetectDelimiter2(t *testing.T) {
	detector := New()

	file, err := os.OpenFile("./Fixtures/test2.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file.Close()

	delimiters := detector.DetectDelimiter(file, '"')
	fmt.Println(delimiters)
	assert.Equal(t, []string{","}, delimiters)
}

func TestDetectRowTerminator(t *testing.T) {
	detector := New()

	file, err := os.OpenFile("./Fixtures/test1.csv", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file.Close()

	terminator := detector.DetectRowTerminator(file)
	assert.Equal(t, "\n", terminator)
}
