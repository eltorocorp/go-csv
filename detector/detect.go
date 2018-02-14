package detector

import (
	"bytes"
	"io"
	"regexp"

	"github.com/jfyne/csvd"
)

const (
	sampleLines             = 15
	nonDelimiterRegexString = `[[:alnum:]\n\r]`
)

// New a detector.
func New() Detector {
	return &detector{
		nonDelimiterRegex: regexp.MustCompile(nonDelimiterRegexString),
	}
}

// Detector defines the exposed interface.
type Detector interface {
	DetectDelimiter(reader io.Reader, enclosure byte) []string
	DetectRowTerminator(reader io.Reader) string
}

// detector is the default implementation of Detector.
type detector struct {
	nonDelimiterRegex *regexp.Regexp
}

// DetectRowTerminator finds the the row terminating string
func (d *detector) DetectRowTerminator(reader io.Reader) string {
	KB := 1024
	buf := make([]byte, 128*KB)
	_, err := reader.Read(buf)
	if err != nil {
		if err == io.EOF {
			return ""
		}
		return ""
	}

	if bytes.Index(buf, []byte{'\r', '\n'}) != -1 {
		return "\r\n"
	}
	if bytes.Index(buf, []byte{'\r'}) != -1 {
		return "\r"
	}

	return "\n"
}

// DetectDelimiter finds a slice of delimiter string.
func (d *detector) DetectDelimiter(r io.Reader, enclosure byte) []string {
	csvReader := csvd.NewReader(r)
	delimiter := csvReader.Comma
	return []string{string(delimiter)}
}
