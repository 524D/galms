package digest

import (
	"reflect"
	"testing"
)

func TestDigestor_Cut(t *testing.T) {
	dSimple := New(0, 0, nil, nil)
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
				t.Errorf("Digestor.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	dUncleaved11 := New(1, 1, nil, nil)
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
				t.Errorf("Digestor.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	dUncleaved01 := New(0, 1, nil, nil)
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
				t.Errorf("Digestor.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
	f5to10 := func(s string) bool { l := len(s); return l >= 5 && l <= 10 }
	dLen5to10 := New(0, 0, f5to10, nil)
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
				t.Errorf("Digestor.Cut() = %v, want %v", got, tt.want)
			}
		})
	}
}
