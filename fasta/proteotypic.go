// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package fasta

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/524D/galms/digest"
)

type proteoTypic struct {
	prot   Prot
	pepInf []pepInfo
}

type pepInfo struct {
	pep   string
	count int
}

// ProteotypicPeps computes proteotypic peptides
// Proteotypic peptides are peptides that are unique for a specific protein
// Unique should be interpreted in a Mass Spectrometric way:
//  equal mass amino acids are treated as the same.
func ProteotypicPeps(fastas []Fasta, enzyme digest.Enzyme) []proteoTypic {
	// Count number of occurrences in different proteins of each peptide
	isoPepCount := make(map[string]int)
	dig := digest.New(0, 0, nil, enzyme)
	// Pass one: determine which peptides are proteotypic
	for _, f := range fastas {
		for _, p := range f.Prots() {
			seq := p.Sequence()
			peps := dig.Cut(seq)
			peps = removeInvalidSeq(peps) // Remove peptides that contain invalid amino acid codes
			// Convert isoleucines to leucines
			for i := range peps {
				peps[i] = isofyPep(peps[i])
			}
			// Only keep uniq results, so that proteotypic peptides that occur
			// multiple times in a protein are still recognized as proteotypic
			peps, _ = unique(peps)
			for _, p := range peps {
				isoPepCount[p]++
			}
		}
	}
	// Pass 2: collect proteotypic peptide details
	result := make([]proteoTypic, 0)
	for _, f := range fastas {
		for _, p := range f.Prots() {
			var pt proteoTypic
			pt.prot = p
			pPeps := make([]string, 0)
			seq := p.Sequence()
			peps := dig.Cut(seq)
			peps = removeInvalidSeq(peps) // Remove peptides that contain invalid amino acid codes

			// If is proteotypic, add it to list
			for _, p := range peps {
				isoPep := isofyPep(p)
				if isoPepCount[isoPep] == 1 {
					pPeps = append(pPeps, p)
				}
			}
			// Check if there are any proteotypic peptides
			if len(pPeps) > 0 {
				pPeps, pepCount := unique(pPeps)
				// Order peptides by length, then alphabetic
				sort.Slice(pPeps, func(i int, j int) bool {
					k := len(pPeps[i]) - len(pPeps[j])
					if k == 0 {
						return pPeps[i] < pPeps[j]
					}
					return k < 0
				})

				for _, p := range pPeps {
					var pi pepInfo
					pi.pep = p
					pi.count = pepCount[p]
					pt.pepInf = append(pt.pepInf, pi)
				}
			}
			result = append(result, pt)
		}
	}
	return result
}

// WriteProteotypicPeps outputs proteotypic peptide info in a human readable format
// First, proteins with proteotypic peptides are written.
// For each protein all proteotypic peptides are listed. Proteotypic peptides that occur multiple
// times in the same protein should be marked with the multiplicity between brackets e.g. (3)
// Next, proteins without proteotypic peptides are written.
func WriteProteotypicPeps(fp *os.File, pts []proteoTypic) {
	// Write proteins with proteotypic peps
	sep := ``
	for _, p := range pts {
		if p.pepInf != nil {
			fmt.Fprintf(fp, "%s%s %s\n", sep, p.prot.ID(), p.prot.Description())
			sep = "\n"

			for _, p := range p.pepInf {
				fmt.Fprint(fp, p.pep)
				if p.count > 1 {
					fmt.Fprintf(fp, " (%d)", p.count)
				}
				fmt.Fprint(fp, "\n")
			}
		}
	}
	// Write proteins without proteotypic peptides
	firstTxt := "\nProteins without proteotypic peptides\n"
	for _, p := range pts {
		if p.pepInf == nil {
			fmt.Fprintf(fp, "%s%s %s\n", firstTxt, p.prot.ID(), p.prot.Description())
			firstTxt = ``
		}
	}
}

// unique finds unique strings
// Returns a slice with unique strings and a maps with occurrence count of each string
func unique(strs []string) ([]string, map[string]int) {
	cnt := make(map[string]int)

	for _, str := range strs {
		cnt[str]++
	}

	ustr := make([]string, 0, len(cnt))
	for p := range cnt {
		ustr = append(ustr, p)
	}
	return ustr, cnt
}

// Map same mass amino acids (I and L) to a single symbol (L)
func isofyPep(pep string) string {
	return strings.Replace(pep, "I", "L", -1)
}

// Remove peptides that contian invalid characters from slice
func removeInvalidSeq(peps []string) []string {
	validPeps := make([]string, 0)
	for _, pep := range peps {
		if !strings.ContainsAny(pep, "BJXZ") {
			validPeps = append(validPeps, pep)
		}
	}
	return validPeps
}
