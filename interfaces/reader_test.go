// Copyright 2014 Jens Rantil. All rights reserved.  Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package interfaces

import (
	"bytes"

	"github.com/eltorocorp/go-csv"

	"testing"

	
)

func TestReaderInterface(t *testing.T) {
	t.Parallel()

	var iface Reader
	iface = csv.NewReader(new(bytes.Buffer))
	iface = csv.NewDialectReader(new(bytes.Buffer), csv.Dialect{})
	iface = csv.NewReader(new(bytes.Buffer))

	// To get rid of compile-time warning that this variable is not used.
	iface.Read()
}
