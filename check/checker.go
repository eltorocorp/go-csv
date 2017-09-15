package check

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

var All = []encoding.Encoding{
    UTF8,
    UTF16(BigEndian, UseBOM),
    UTF16(BigEndian, IgnoreBOM),
    UTF16(LittleEndian, IgnoreBOM),
}


