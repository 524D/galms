// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package molecule

import (
	"reflect"
	"testing"

	"github.com/524D/galms/elements"
)

func TestSimpleFormula(t *testing.T) {
	type args struct {
		f string
		e *elements.Elems
	}
	elms := elements.New()
	tests := []struct {
		name    string
		args    args
		want    Molecule
		wantErr bool
	}{
		{
			name: "Correct formula1",
			args: args{`H2SO4`, elms},
			want: Molecule{
				atoms: []AtomsCount{
					{idx: 0, count: 2},
					{idx: 7, count: 4},
					{idx: 15, count: 1},
				},
				e: elms,
			},
			wantErr: false,
		},
		{
			name: "Correct formula2",
			args: args{`NaCl`, elms},
			want: Molecule{
				atoms: []AtomsCount{
					{idx: 10, count: 1},
					{idx: 16, count: 1},
				},
				e: elms,
			},
			wantErr: false,
		},
		{
			name: "Re occurring elements",
			args: args{`H2C3OH`, elms},
			want: Molecule{
				atoms: []AtomsCount{
					{idx: 0, count: 3},
					{idx: 5, count: 3},
					{idx: 7, count: 1},
				},
				e: elms,
			},
			wantErr: false,
		},
		{
			name:    "Incorrect formula1",
			args:    args{`NaCw`, elms},
			want:    Molecule{},
			wantErr: true,
		},
		{
			name:    "Incorrect formula2",
			args:    args{`Na3333333333333333333333333333`, elms},
			want:    Molecule{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleFormula(tt.args.f, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleFormula() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleFormula() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initAA(t *testing.T) {
	type args struct {
		e *elements.Elems
	}
	elms := elements.New()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Init test",
			args: args{elms},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCommonMolecules(tt.args.e)
		})
	}
}

func TestChemicalFormula(t *testing.T) {
	type args struct {
		m Molecule
		e *elements.Elems
	}
	elms := elements.New()
	fh2o, _ := SimpleFormula("H2O", elms)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Water test",
			args:    args{fh2o, elms},
			want:    "H2O",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChemicalFormula(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChemicalFormula() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChemicalFormula() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		m1 Molecule
		m2 Molecule
	}
	elms := elements.New()
	fh2o, _ := SimpleFormula("H2O", elms)
	fcs, _ := SimpleFormula("CS", elms)
	fo, _ := SimpleFormula("O", elms)
	fh2cos, _ := SimpleFormula("H2COS", elms)
	fh2o2, _ := SimpleFormula("H2O2", elms)
	tests := []struct {
		name string
		args args
		want Molecule
	}{
		{
			name: "H2O + CS",
			args: args{fh2o, fcs},
			want: fh2cos,
		},
		{
			name: "H2O + O",
			args: args{fh2o, fo},
			want: fh2o2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPepProt(t *testing.T) {
	type args struct {
		p string
	}
	elms := elements.New()
	fcystine, _ := SimpleFormula("C3H5NOSH2O", elms)

	tests := []struct {
		name    string
		args    args
		want    Molecule
		wantErr bool
	}{
		{
			name:    "Cystine",
			args:    args{`C`},
			want:    fcystine,
			wantErr: false,
		},
		{
			name:    "Non existing AA",
			args:    args{`CB`},
			want:    Molecule{atoms: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PepProt(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("PepProt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PepProt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAminoAcid(t *testing.T) {
	type args struct {
		aa byte
	}
	elms := elements.New()
	fcystinenw, _ := SimpleFormula("C3H5NOS", elms)
	tests := []struct {
		name    string
		args    args
		want    Molecule
		wantErr bool
	}{
		{
			name:    "Cystine",
			args:    args{'C'},
			want:    fcystinenw,
			wantErr: false,
		},
		{
			name:    "Non existing AA",
			args:    args{'B'},
			want:    Molecule{atoms: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AminoAcid(tt.args.aa)
			if (err != nil) != tt.wantErr {
				t.Errorf("AminoAcid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AminoAcid() = %v, want %v", got, tt.want)
			}
		})
	}
}
