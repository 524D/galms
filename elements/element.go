package elements

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
)

// Element contains mass spec related information about an chemical element
type Element struct {
	Symbol  string    // Symbol, e.g. "Na"
	Name    string    // Full name, e.g. sodium
	Number  int       // Atomic number e.g. 11
	Isotope []Isotope // Isotope masses and abundances, e.g. {22.989770,1}
}

// Isotope contains mass and abundance of a single isotope
type Isotope struct {
	Mass      float64 // Mass in Daltons
	Abundance float64 // Abundance (fraction)
}

// Elems wraps the contents of all elements
type Elems struct {
	Elements  []Element
	SymbolMap map[string]int
}

// Read reads elements in JSON format from an io.Reader
func Read(r io.Reader) ([]Element, error) {
	var elements []Element

	err := json.NewDecoder(r).Decode(&elements)
	return elements, err
}

func Read1(r io.Reader) (Elems, error) {
	var e Elems

	e.SymbolMap = make(map[string]int)
	err := json.NewDecoder(r).Decode(&e.Elements)
	if err != nil {
		return e, err
	}
	for i, ei := range e.Elements {
		e.SymbolMap[ei.Symbol] = i
	}

	return e, err
}

// InitDefault parses the build-in list of elements
func InitDefault() Elems {
	e, err := Read1(strings.NewReader(defaultElementsJSON))
	if err != nil {
		// Should never happen
		log.Fatal("Error reading build-in elements")
	}
	return e
}

// ElemIdx converts an element string to an index
func (e *Elems) ElemIdx(shortName string) (int, error) {
	i, ok := e.SymbolMap[shortName]
	if !ok {
		return 0, errors.New("ElemIdx: Element not found")
	}
	return i, nil
}

// Name returns the element name for a given element index
func (e *Elems) Name(i int) (string, error) {
	if i < 0 && i >= len(e.Elements) {
		return ``, errors.New("Name: Index out of range")
	}
	return e.Elements[i].Name, nil
}

// Symbol returns the element name for a given element index
func (e *Elems) Symbol(i int) (string, error) {
	if i < 0 && i >= len(e.Elements) {
		return ``, errors.New("Symbol: Index out of range")
	}
	return e.Elements[i].Symbol, nil
}

// Isotopes retruns the isotopes of a specific element
func (e *Elems) Isotopes(i int) ([]Isotope, error) {
	if i < 0 && i >= len(e.Elements) {
		return nil, errors.New("Isotopes: Index out of range")
	}
	return e.Elements[i].Isotope, nil
}
