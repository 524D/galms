package molecule

import (
	"errors"
	"regexp"
	"sort"
	"strconv"

	"github.com/524D/galms/elements"
)

// This package parses chemical formula

// AtomsCount contains an atom index and atom count
type AtomsCount struct {
	idx   int // Element index
	count int // Number of atoms of this element
}

// Molecule represents a single molecule
type Molecule struct {
	atoms []AtomsCount
}

// ErrUnknownAACode Single letter AA code unknown
var ErrUnknownAACode = errors.New("Unknown amino acid code")

// Atoms returns the atoms/count of a molecule
func (m *Molecule) Atoms() []AtomsCount {
	return m.atoms
}

// IdxCount retruns the index and count of an atom in a molecule
func (ac *AtomsCount) IdxCount() (int, int) {
	return ac.idx, ac.count
}

// SimpleFormula converts a chemical formula to a structure that contains elements/atom counts
func SimpleFormula(f string, e *elements.Elems) (Molecule, error) {
	var m Molecule

	// Map previously encountered element index into position in 'atoms'
	prevEl := make(map[int]int)
	re := regexp.MustCompile(`([A-Z][a-z]?)([0-9]*)`)
	mts := re.FindAllStringSubmatch(f, -1)
	for _, em := range mts {
		var elCount AtomsCount
		var err error
		// em[1] holds the elements string
		elCount.idx, err = e.ElemIdx(em[1])
		if err != nil {
			return m, err
		}
		if em[2] == `` {
			elCount.count = 1
		} else {
			elCount.count, err = strconv.Atoi(em[2])
			if err != nil {
				return m, err
			}
		}
		// If we already encountered this element previously,
		// just add the atom count
		if i, ok := prevEl[elCount.idx]; ok {
			m.atoms[i].count += elCount.count
		} else {
			prevEl[elCount.idx] = len(m.atoms)
			m.atoms = append(m.atoms, elCount)
		}
	}
	// Sort atoms in element index
	sort.Slice(m.atoms, func(i, j int) bool { return m.atoms[i].idx < m.atoms[j].idx })
	return m, nil
}

// ChemicalFormula a molecule to a string
func ChemicalFormula(m Molecule, e *elements.Elems) (string, error) {
	var f string

	for _, a := range m.atoms {
		// Ignore elements with zero (or negative) count
		if a.count > 0 {
			symbol, err := e.Symbol(a.idx)
			if err != nil {
				return ``, err
			}
			f += symbol
			if a.count > 1 {
				f += strconv.Itoa(a.count)
			}
		}
	}
	return f, nil
}

// AminoAcid returns the molecule (minus H2O) for a single letter amino acid code
func AminoAcid(aa byte) (Molecule, error) {
	var m Molecule
	if aaMol[aa].atoms == nil {
		return m, ErrUnknownAACode
	}
	return aaMol[aa], nil
}

// PepProt returns the molecule peptide or protein
func PepProt(p string) (Molecule, error) {
	var m Molecule

	for _, aa := range p {
		if aa >= 'A' && aa <= 'Z' && aaMol[aa].atoms != nil {
			m = Add(m, aaMol[aa])
		} else {
			return Molecule{}, ErrUnknownAACode
		}
	}
	m = Add(m, water)

	return m, nil
}

// Add two molecules. The atoms in the molecules must be sorted by atom index
func Add(m1 Molecule, m2 Molecule) Molecule {
	var m Molecule
	// Make sure 2 have enough room even if both molecules contain
	// completely different atoms
	m.atoms = make([]AtomsCount, len(m1.atoms), len(m1.atoms)+len(m2.atoms))

	copy(m.atoms, m1.atoms)
	mi := 0

	for i, a := range m2.atoms {
		for mi < len(m.atoms) && m.atoms[mi].idx < a.idx {
			mi++
		}
		// If we are at the end of the first molecules atoms, append the rest of the
		// second one, and we are done
		if mi >= len(m.atoms) {
			m.atoms = append(m.atoms, m2.atoms[i:]...)
			return m
		}
		// If same atom index, add the count of this atom
		if m.atoms[mi].idx == a.idx {
			m.atoms[mi].count += a.count
		} else {
			// m2 contains an atom index that was not in m, insert the atom
			// This is 'expensive', but normally doesn't happen often
			// Extend m.atoms by one (last element is overwritten later)
			m.atoms = append(m.atoms, a)
			copy(m.atoms[mi+1:], m.atoms[mi:])
			m.atoms[mi] = a
		}
	}

	return m
}

// Conversion table for translating amino acids to molecules
// Initialized by 'initAA'
var aaMol [256]Molecule
var water Molecule

type aaForm struct {
	code    byte
	formula string
}

// InitCommonMolecules initializes some common molecules
// It must be called after (re-)initializing elements.Elems
func InitCommonMolecules(e *elements.Elems) {
	// Set up amino acid translation table
	aaList := []aaForm{
		{code: 'A', formula: `C3H5NO`},
		{code: 'C', formula: `C3H5NOS`},
		{code: 'D', formula: `C4H5NO3`},
		{code: 'E', formula: `C5H7NO3`},
		{code: 'F', formula: `C9H9NO`},
		{code: 'G', formula: `C2H3NO`},
		{code: 'H', formula: `C6H7N3O`},
		{code: 'I', formula: `C6H11NO`},
		{code: 'K', formula: `C6H12N2O`},
		{code: 'L', formula: `C6H11NO`},
		{code: 'M', formula: `C5H9NOS`},
		{code: 'N', formula: `C4H6N2O2`},
		{code: 'O', formula: `C5H7NO2`},
		{code: 'P', formula: `C5H7NO`},
		{code: 'Q', formula: `C5H8N2O2`},
		{code: 'R', formula: `C6H12N4O`},
		{code: 'S', formula: `C3H5NO2`},
		{code: 'T', formula: `C4H7NO2`},
		{code: 'U', formula: `C5H5NO2`},
		{code: 'V', formula: `C5H9NO`},
		{code: 'W', formula: `C11H10N2O`},
		{code: 'Y', formula: `C9H9NO2`},
	}
	for i := range aaMol {
		aaMol[i].atoms = nil
	}
	for _, a := range aaList {
		aaMol[a.code], _ = SimpleFormula(a.formula, e)
	}
	water, _ = SimpleFormula(`H2O`, e)
}
