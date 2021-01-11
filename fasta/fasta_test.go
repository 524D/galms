package fasta

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Fasta
		wantErr bool
	}{
		{
			name: "Simple FASTA",
			args: args{reader: strings.NewReader(`>BLABLA Something
ACDEFGH
>BLA2
EEEY
>
WE
ARE
CPM
`)},
			want: Fasta{
				[]Prot{
					{ID: "BLABLA", desc: "Something", seq: "ACDEFGH"},
					{ID: "BLA2", desc: "", seq: "EEEY"},
					{ID: "DUMMY_ID_3", desc: "", seq: "WEARECPM"},
				},
			},
			wantErr: false,
		},
		{
			name: "No newline end FASTA",
			args: args{reader: strings.NewReader(`>TEST2 No newline at end
AHAH`)},
			want: Fasta{
				[]Prot{
					{ID: "TEST2", desc: "No newline at end", seq: "AHAH"},
				},
			},
			wantErr: false,
		},
		{
			name: "Additional spacing FASTA",
			args: args{reader: strings.NewReader(`>TEST3  	 After some spaces and tab
HAHA
`)},
			want: Fasta{
				[]Prot{
					{ID: "TEST3", desc: "After some spaces and tab", seq: "HAHA"},
				},
			},
			wantErr: false,
		},
		{
			name: "Additional spacing FASTA2",
			args: args{reader: strings.NewReader(`>TEST4 Spaces in/around seq
  HADIHI 

  NAH 

`)},
			want: Fasta{
				[]Prot{
					{ID: "TEST4", desc: "Spaces in/around seq", seq: "HADIHINAH"},
				},
			},
			wantErr: false,
		},

		// FIXME: Add test that generates error
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
