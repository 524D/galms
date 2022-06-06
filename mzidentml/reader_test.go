// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package mzidentml

import (
	"log"
	"os"
	"testing"
)

// Test files, downloaded Pride
// const testFile1 = "/home/robm/data/mzid_testfiles/F181955 -filtered.mzid"
// const testFile2 = "/home/robm/data/mzid_testfiles/BalbfedbyRag_Day7_01.mzid"

// Test files, referenced from mzID specification at https://www.psidev.info/mzidentml
// downloaded from https://github.com/HUPO-PSI/mzIdentML/tree/master/examples/
const testFile1 = "./testfiles/55merge_tandem.mzid"
const testFile2 = "./testfiles/combined_1.2.mzid"

func TestAll1(t *testing.T) {

	x, err := os.Open(testFile1)
	if err != nil {
		t.Errorf("Open: mzIdentMLfile is nil")
	}
	f, err := Read(x)
	if err != nil {
		t.Errorf("Read: error return %v", err)
	}
	//	log.Printf("%+v\n", f.content)

	n := f.NumIdents()
	if n != 170 {
		t.Errorf("NumIdents is %d, expected 170", n)
	}
	ident, err := f.Ident(100)
	if err != nil {
		t.Errorf("Ident: error return %v", err)
	}
	// log.Printf("ident: %+v", ident)
	if ident.PepSeq != `FLLSEVGPMSAR` {
		t.Errorf("Seqence 100=%s, expected FLLSEVGPMSAR", ident.PepSeq)
	}
}

func TestAll2(t *testing.T) {

	x, err := os.Open(testFile2)
	if err != nil {
		t.Errorf("Open: mzIdentMLfile is nil")
	}
	f, err := Read(x)
	if err != nil {
		t.Errorf("Read: error return %v", err)
	}
	//	log.Printf("%+v\n", f.content)

	n := f.NumIdents()
	if n != 11707 {
		t.Errorf("NumIdents is %d, expected 8754", n)
	}
	ident, err := f.Ident(1000)
	if err != nil {
		t.Errorf("Ident: error return %v", err)
	}
	log.Printf("ident: %+v", ident)

}
