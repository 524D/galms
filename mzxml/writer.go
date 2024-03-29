// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package mzxml

import (
	"encoding/xml"
	"io"
)

func (mzXML *MzXML) Write(w io.Writer) error {
	content := &mzXML.content
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "    ")
	err := encoder.Encode(content)
	return err
}
