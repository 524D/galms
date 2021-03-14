package digest

// Enzyme must return "true" if enzyme cuts at position pos in sequence seq
type Enzyme func(seq string, pos int) bool

// Enzyme must return "true" if enzyme fullfils creteria (e.g. length)
type Filter func(seq string) bool

type Digestor struct {
	minMissedCleavage int
	maxMissedCleavage int
	filter            Filter
	enzyme            Enzyme
}

func New(minMissedCleavage int, maxMissedCleavage int, filter Filter, enzyme Enzyme) *Digestor {
	var d Digestor

	d.filter = filter
	d.minMissedCleavage = minMissedCleavage
	d.maxMissedCleavage = maxMissedCleavage
	if enzyme == nil {
		d.enzyme = Trypsin
	} else {
		d.enzyme = enzyme
	}
	return &d
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Trypsin(seq string, i int) bool {
	// Never happens:
	// if i < 1 {
	// 	return false
	// }
	c1 := seq[i-1]
	c2 := seq[i]
	return (c1 == 'K' || c1 == 'R') && c2 != 'P'
}

func (d *Digestor) cleave(seq string) []string {
	p := make([]string, 0, 20)
	prev := 0
	for i := 1; i < len(seq)-1; i++ {
		if d.enzyme(seq, i) {
			p = append(p, seq[prev:i])
			prev = i
		}
	}
	p = append(p, seq[prev:])
	return p
}

func (d *Digestor) Cut(seq string) []string {
	peps := make([]string, 0, 20)
	p := d.cleave(seq)
	// To compose a list of all peptides with missed cleavages,
	for skip := 0; skip <= d.maxMissedCleavage; skip++ {
		glueMin := maxInt(d.minMissedCleavage, skip)
		for glue := glueMin; glue <= d.maxMissedCleavage; glue++ {
			for i := skip; i < len(p)-glue; i += glue + 1 {
				pep := p[i]
				for j := 1; j <= glue; j++ {
					pep += p[i+j]
				}
				if d.filter == nil || d.filter(pep) {
					peps = append(peps, pep)
				}
			}
		}
	}
	return peps
}
