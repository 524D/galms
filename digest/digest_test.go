// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package digest

import (
	"reflect"
	"runtime"
	"testing"
)

func TestDigestor_Cut(t *testing.T) {
	dSimple := New(0, 0, nil, TrypsinSimple)
	type args struct {
		seq string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Single cleave",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`MK`, `WVTFISLLFLFSSAYSR`, `GVFR`, `R`, `DAHK`, `SEVAHR`, `FK`, `DLGEENFK`, `ALVLIAFAQYLQQCPF`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dSimple.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	dUncleaved11 := New(1, 1, nil, TrypsinSimple)
	tests2 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Missed cleavage",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`MKWVTFISLLFLFSSAYSR`, `GVFRR`, `DAHKSEVAHR`, `FKDLGEENFK`, `WVTFISLLFLFSSAYSRGVFR`, `RDAHK`, `SEVAHRFK`, `DLGEENFKALVLIAFAQYLQQCPF`},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := dUncleaved11.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	dUncleaved01 := New(0, 1, nil, TrypsinSimple)
	tests3 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Missed cleavage",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`MK`, `WVTFISLLFLFSSAYSR`, `GVFR`, `R`, `DAHK`, `SEVAHR`, `FK`, `DLGEENFK`, `ALVLIAFAQYLQQCPF`, `MKWVTFISLLFLFSSAYSR`, `GVFRR`, `DAHKSEVAHR`, `FKDLGEENFK`, `WVTFISLLFLFSSAYSRGVFR`, `RDAHK`, `SEVAHRFK`, `DLGEENFKALVLIAFAQYLQQCPF`},
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			if got := dUncleaved01.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	f5to10 := func(s string) bool { l := len(s); return l >= 5 && l <= 10 }
	dLen5to10 := New(0, 0, f5to10, TrypsinSimple)
	tests4 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Single cleave",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`SEVAHR`, `DLGEENFK`},
		},
	}
	for _, tt := range tests4 {
		t.Run(tt.name, func(t *testing.T) {
			if got := dLen5to10.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	rCutter := func(seq string, pos int) bool {
		return (seq[pos-1] == 'R')
	}
	dRCutter := New(0, 0, nil, rCutter)
	tests5 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Special enzyme",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`MKWVTFISLLFLFSSAYSR`, `GVFR`, `R`, `DAHKSEVAHR`, `FKDLGEENFKALVLIAFAQYLQQCPF`},
		},
	}
	for _, tt := range tests5 {
		t.Run(tt.name, func(t *testing.T) {
			if got := dRCutter.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}

	e, _ := NamedEnzyme(`chymotrypsin`)
	cutter := New(0, 0, nil, e)
	tests6 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Chymotrypsin",
			args: args{seq: `MKWVTFISLLFLFSSAYSRGVFRRDAHKSEVAHRFKDLGEENFKALVLIAFAQYLQQCPF`},
			want: []string{`MKW`, `VTF`, `ISL`, `L`, `F`, `L`, `F`, `SSAY`, `SRGVF`, `RRDAHKSEVAHRF`, `KDL`, `GEENF`, `KAL`, `VL`, `IAF`, `AQY`, `L`, `QQCPF`},
		},
	}
	for _, tt := range tests6 {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutter.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}

	cutter = New(0, 0, nil, nil)
	tests7 := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Real Trypsin, with all special cases",
			args: args{seq: `RRRMARAAKGGRPWKPEERPEMRPAACKDADKDRAACKHACKYAAKYCRKAARRHAARRAAA`},
			want: []string{`R`, `RR`, `MAR`, `AAK`, `GGRPWK`, `PEERPEMR`, `PAACKDADKDR`, `AACKHACKYAAK`, `YCRK`, `AAR`, `RHAAR`, `R`, `AAA`},
		},
	}
	for _, tt := range tests7 {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutter.Cut(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("digest.Cut() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestEnzymes(t *testing.T) {
	tests := []struct {
		name string
		want []EnzymeInf
	}{
		{
			name: "",
			want: []EnzymeInf{
				{`Trypsin`, `See https://web.expasy.org/peptide_cutter/peptidecutter_enzymes.html`, Trypsin},
				{`Trypsin_Simple`, `Cuts after K and R but not before P`, TrypsinSimple},
				{`Trypsin/P`, `Cuts after K and R`, TrypsinP},
				{`Lys_C`, `Cuts after K but not before P`, LysC},
				{`PepsinA`, `Cuts after F and L but not before P`, PepsinA},
				{`Chymotrypsin`, `Cuts after F, W, Y, L but not before P`, ChymoTrypsin},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Can't use Deepequal to compare function pointers, so we do it manually
			got := Enzymes()
			for i := range tt.want {
				name1 := tt.want[i].Name
				name2 := got[i].Name
				desc1 := tt.want[i].Description
				desc2 := got[i].Description
				funcName1 := runtime.FuncForPC(reflect.ValueOf(tt.want[i].Func).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(got[i].Func).Pointer()).Name()
				if name1 != name2 || desc1 != desc2 || funcName1 != funcName2 {
					t.Errorf("Enzymes() = %v, want %v", got, tt.want)

				}
			}
		})
	}
}

func TestNamedEnzyme(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		args    args
		want    Enzyme
		wantErr bool
	}{
		{
			name:    "Trypsin test",
			args:    args{`trypsin`},
			want:    Trypsin,
			wantErr: false,
		},
		{
			name:    "Invalid enzyme name test",
			args:    args{`trapsin`},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamedEnzyme(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamedEnzyme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			funcName1 := runtime.FuncForPC(reflect.ValueOf(tt.want).Pointer()).Name()
			funcName2 := runtime.FuncForPC(reflect.ValueOf(got).Pointer()).Name()
			if funcName1 != funcName2 {
				t.Errorf("NamedEnzyme() = %v, want %v", got, tt.want)

			}
		})
	}
}
