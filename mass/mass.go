package mass

import (
	"math"

	"../elements"
	"../molecule"
)

// Peak contains the mass and abundance of one of the isotopes combinations of a molecule
type Peak struct {
	Mass      float64
	Abundance float64
}

// MinMax returns the smallest and largest mass of a molecule
// The relative abundance is also returned
func MinMax(m molecule.Molecule, e *elements.Elems) (Peak, Peak, error) {
	min := Peak{Mass: 0, Abundance: 1.0}
	max := Peak{Mass: 0, Abundance: 1.0}

	for _, a := range m.Atoms() {
		idx, count := a.IdxCount()
		iso, err := e.Isotopes(idx)
		if err != nil {
			return Peak{}, Peak{}, nil
		}
		min.Mass += iso[0].Mass * float64(count)
		min.Abundance *= math.Pow(iso[0].Abundance, float64(count))
		max.Mass += iso[len(iso)-1].Mass * float64(count)
		min.Abundance *= math.Pow(iso[len(iso)-1].Abundance, float64(count))
	}
	return min, max, nil
}
