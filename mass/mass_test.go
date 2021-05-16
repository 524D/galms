// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package mass

import (
	"math"
	"testing"

	"github.com/524D/galms/elements"
	"github.com/524D/galms/molecule"
)

const tolerance = float64(0.00001)

func float64Eql(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return true
	}
	return (diff / mean) < tolerance
}

func TestMinMax(t *testing.T) {
	type args struct {
		m molecule.Molecule
		e *elements.Elems
	}
	elms := elements.New()
	mol1, _ := molecule.SimpleFormula("C10", elms)
	tests := []struct {
		name    string
		args    args
		want    Peak
		want1   Peak
		wantErr bool
	}{
		{
			name: "Mass test1",
			args: args{
				m: mol1,
				e: elms,
			},
			want: Peak{
				Mass:      120.0,
				Abundance: 0.898007762,
			},
			want1: Peak{
				Mass:      130.0335483507,
				Abundance: 1.967151357e-20,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := MinMax(tt.args.m, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinMax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !float64Eql(got.Abundance, tt.want.Abundance) ||
				!float64Eql(got.Mass, tt.want.Mass) {
				t.Errorf("MinMax() got = %v, want %v", got, tt.want)
			}
			if !float64Eql(got1.Abundance, tt.want1.Abundance) ||
				!float64Eql(got1.Mass, tt.want1.Mass) {
				t.Errorf("MinMax() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}
