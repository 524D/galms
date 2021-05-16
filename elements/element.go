// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package elements

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
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

var once sync.Once
var defaultElms *Elems

// Read reads elements in JSON format from an io.Reader
func Read(r io.Reader) (*Elems, error) {
	var e Elems

	err := json.NewDecoder(r).Decode(&e.Elements)
	if err != nil {
		return nil, err
	}
	e.SymbolMap = make(map[string]int)
	for i, ei := range e.Elements {
		e.SymbolMap[ei.Symbol] = i
	}

	return &e, err
}

// New parses the build-in list of elements
func New() *Elems {
	once.Do(func() {
		defaultElms, _ = Read(strings.NewReader(defaultElementsJSON))
	})
	return defaultElms
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
	if i < 0 || i >= len(e.Elements) {
		return ``, errors.New("Name: Index out of range")
	}
	return e.Elements[i].Name, nil
}

// Symbol returns the element symbol for a given element index
func (e *Elems) Symbol(i int) (string, error) {
	if i < 0 || i >= len(e.Elements) {
		return ``, errors.New("Symbol: Index out of range")
	}
	return e.Elements[i].Symbol, nil
}

// Isotopes returns the isotopes of a specific element
func (e *Elems) Isotopes(i int) ([]Isotope, error) {
	if i < 0 || i >= len(e.Elements) {
		return nil, errors.New("Isotopes: Index out of range")
	}
	return e.Elements[i].Isotope, nil
}
