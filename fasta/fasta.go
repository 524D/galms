package fasta

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Fasta wraps the contents of the FASTA file
type Fasta struct {
	prot []Prot
}

// Prot single entry (identifier, description and sequence)
type Prot struct {
	ID   string
	desc string
	seq  string
}

func appendProt(prot Prot, fasta *Fasta) {
	if prot.ID != `` || prot.desc != `` || prot.seq != `` {
		fasta.prot = append(fasta.prot, prot)
	}
}

// Read reads an FASTA file from an io.Reader
func Read(reader io.Reader) (Fasta, error) {
	var fasta Fasta
	var prot Prot
	re := regexp.MustCompile(`>([^ \t]*)(?:[ \t]+(.*))?`)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, ">") {
			appendProt(prot, &fasta)
			m := re.FindStringSubmatch(l)
			if m == nil || len(m) < 2 || m[1] == `` {
				// Parsing ID/description failed
				// Make up a fake ID
				prot.ID = `DUMMY_ID_` + strconv.FormatInt(int64(len(fasta.prot))+1, 10)
			} else {
				prot.ID = m[1]
				if len(m) >= 3 {
					prot.desc = m[2]
				} else {
					prot.desc = ``
				}
			}
			prot.seq = ``
		} else {
			// Add to sequence, remove superfluous spacing
			prot.seq += strings.TrimSpace(l)
		}
	}
	appendProt(prot, &fasta)

	err := scanner.Err()

	return fasta, err
}

// Write writes a new FASTA file to an io.writer
func (f *Fasta) Write(writer io.Writer) error {
	return nil
}
