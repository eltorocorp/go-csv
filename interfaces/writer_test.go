// Copyright 2014 Jens Rantil. All rights reserved.  Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package interfaces

import (
	"bytes"

	"encoding/csv"

	"github.com/eltorocorp/go-csv"

	"testing"

	
)

func TestWriterInterface(t *testing.T) {
	t.Parallel()

	var iface Writer
	iface = csv.NewWriter(new(bytes.Buffer))
	iface = csv.NewDialectWriter(new(bytes.Buffer), csv.Dialect{})
	iface = csv.NewWriter(new(bytes.Buffer))

	// To get rid of compile-time warning that this variable is not used.
	iface.Flush()
}
