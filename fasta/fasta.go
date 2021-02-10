package fasta

import (
	"bufio"
	"fmt"
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
	id   string
	desc string
	seq  string
}

func (f *Fasta) appendProt(prot Prot) {
	if prot.id != `` || prot.desc != `` || prot.seq != `` {
		f.prot = append(f.prot, prot)
	}
}

// Read reads an FASTA file from an io.Reader
func Read(reader io.Reader) (Fasta, error) {
	var fasta Fasta
	var prot Prot
	re := regexp.MustCompile(`>([^ \t]*)(?:[ \t]+(.+)?)?`)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			// Skip PEFF header
		} else if strings.HasPrefix(l, ">") {
			fasta.appendProt(prot)
			m := re.FindStringSubmatch(l)
			if m == nil || len(m) < 2 || m[1] == `` {
				// Parsing ID/description failed
				// Make up a fake ID
				prot.id = `DUMMY_ID_` + strconv.FormatInt(int64(len(fasta.prot))+1, 10)
			} else {
				prot.id = m[1]
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
	fasta.appendProt(prot)

	err := scanner.Err()

	return fasta, err
}

// Write writes a new FASTA file to an io.writer
func (f *Fasta) Write(writer io.Writer) error {
	seqLineLen := 60
	for _, p := range f.prot {
		fmt.Fprintf(writer, ">%s\t%s\n", p.id, p.desc)
		i := 0
		l := ``
		for _, a := range p.seq {
			l += string(a)
			i++
			if i >= seqLineLen {
				_, err := fmt.Fprintf(writer, "%s\n", l)
				if err != nil {
					return err
				}
				i = 0
				l = ``
			}
		}
		if l != `` {
			_, err := fmt.Fprintf(writer, "%s\n", l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Prots returns a slice with all proteins in a fasta file
func (f *Fasta) Prots() []Prot {
	return f.prot
}

// ID returns the protein ID
func (p *Prot) ID() string {
	return p.id
}

// Description returns the protein description
func (p *Prot) Description() string {
	return p.desc
}

// Sequence returns the proten sequence
func (p *Prot) Sequence() string {
	return p.seq
}

// Translate converts genetic sequences into protain sequences
// https://www.ncbi.nlm.nih.gov/Taxonomy/taxonomyhome.html/index.cgi?chapter=cgencodes#SG11
