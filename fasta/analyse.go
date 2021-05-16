package fasta

import (
	"fmt"
	"log"

	"github.com/524D/galms/digest"
	"github.com/524D/galms/elements"
	"github.com/524D/galms/mass"
	"github.com/524D/galms/molecule"
)

// Analyse obtains some common properties for protines in a FASTA file
// - length/mass distribution distribution
// - number of proteins in certain class,
//   e.g.:
//     with motif for export from cell, import golgi, zinc finger, multiple motifs in one prot
//       zinc-binding motif HEXXHXXGXXH
//       N-terminal secretory signal peptide and a prodomain with a conserved PRCGXPD motif
//     isoforms
//     phosphorilation/glycosylation patterns
// - statistically significant sequence features, e.g. missing subsequence/regex, abundant subsequence/regex
// - digest (selectable, default tryptic)
//   for each peptides, print distribution of number of occurrences (unique)
//   print prots with no unique peptides (merging I and L)
//   print prots with no peptides that can be measured or uniquely identified, e.g. because of mass range, ionizability
//
func Analyse(f Fasta, enzyme digest.Enzyme) {
	e := elements.New()
	molecule.InitCommonMolecules(e)
	pepProteins := make(map[string][]Prot)
	maxOccur := 0
	f6to30 := func(s string) bool { l := len(s); return l >= 6 && l <= 30 }
	dig := digest.New(0, 1, f6to30, enzyme)
	sep := ``
	for _, p := range f.Prots() {
		seq := p.Sequence()
		m, err := molecule.PepProt(seq)
		if err != nil {
			log.Printf("Can't convert protein seq to chemical formula for %s %s: %v\n", p.ID(), seq, err)
			continue
		}
		minm, maxm, err := mass.MinMax(m, e)
		if err != nil {
			log.Printf("Can't compute mass for %v: %v\n", m, err)
			continue
		}
		fmt.Printf("%s%s %s\n", sep, p.ID(), p.Description())
		sep = "\n"
		fmt.Printf("Mass min: %f (%f%%) max %f (%f%%)\n", minm.Mass, minm.Abundance, maxm.Mass, maxm.Abundance)
		peps := dig.Cut(seq)
		fmt.Printf("Num peps: %d\n", len(peps))
		for _, pep := range peps {
			pepProteins[pep] = append(pepProteins[pep], p)
			if maxOccur < len(pepProteins[pep]) {
				maxOccur = len(pepProteins[pep])
			}
		}
	}
	occurCnt := make([]string, maxOccur+1)
	for pep, prots := range pepProteins {
		occurCnt[len(prots)] = pep
	}
	fmt.Printf("The most peptide (%s) occurs in %d proteins\n", occurCnt[len(occurCnt)-1], len(occurCnt))

}
