package fasta

import (
	"bytes"
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
					{id: "BLABLA", desc: "Something", seq: "ACDEFGH"},
					{id: "BLA2", desc: "", seq: "EEEY"},
					{id: "DUMMY_ID_3", desc: "", seq: "WEARECPM"},
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
					{id: "TEST2", desc: "No newline at end", seq: "AHAH"},
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
					{id: "TEST3", desc: "After some spaces and tab", seq: "HAHA"},
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
					{id: "TEST4", desc: "Spaces in/around seq", seq: "HADIHINAH"},
				},
			},
			wantErr: false,
		},
		{
			name: "PEFF header",
			args: args{reader: strings.NewReader(`
# This ain't no PEFF file!
>TEST5 Blabla
HAHAHA
`)},
			want: Fasta{
				[]Prot{
					{id: "TEST5", desc: "Blabla", seq: "HAHAHA"},
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

func TestFasta_Prots(t *testing.T) {
	type fields struct {
		prot []Prot
	}
	tests := []struct {
		name   string
		fields fields
		want   []Prot
	}{
		{
			name: "Simple FASTA prots",
			fields: fields{
				prot: []Prot{
					{id: "BLABLA", desc: "Something", seq: "ACDEFGH"},
					{id: "BLA2", desc: "", seq: "EEEY"},
					{id: "DUMMY_ID_3", desc: "", seq: "WEARECPM"},
				},
			},

			want: []Prot{
				{id: "BLABLA", desc: "Something", seq: "ACDEFGH"},
				{id: "BLA2", desc: "", seq: "EEEY"},
				{id: "DUMMY_ID_3", desc: "", seq: "WEARECPM"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fasta{
				prot: tt.fields.prot,
			}
			if got := f.Prots(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fasta.Prots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProt_ID(t *testing.T) {
	type fields struct {
		id   string
		desc string
		seq  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				id:   "SOMEPROT",
				desc: "Cysteine rich peptide",
				seq:  "CCCCCCCCCCCCC",
			},
			want: "SOMEPROT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prot{
				id:   tt.fields.id,
				desc: tt.fields.desc,
				seq:  tt.fields.seq,
			}
			if got := p.ID(); got != tt.want {
				t.Errorf("Prot.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProt_Description(t *testing.T) {
	type fields struct {
		id   string
		desc string
		seq  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				id:   "SOMEPROT",
				desc: "Cysteine rich peptide",
				seq:  "CCCCCCCCCCCCC",
			},
			want: "Cysteine rich peptide",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prot{
				id:   tt.fields.id,
				desc: tt.fields.desc,
				seq:  tt.fields.seq,
			}
			if got := p.Description(); got != tt.want {
				t.Errorf("Prot.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProt_Sequence(t *testing.T) {
	type fields struct {
		id   string
		desc string
		seq  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				id:   "SOMEPROT",
				desc: "Cysteine rich peptide",
				seq:  "CCCCCCCCCCCCC",
			},
			want: "CCCCCCCCCCCCC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prot{
				id:   tt.fields.id,
				desc: tt.fields.desc,
				seq:  tt.fields.seq,
			}
			if got := p.Sequence(); got != tt.want {
				t.Errorf("Prot.Sequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFasta_Write(t *testing.T) {
	type fields struct {
		prot []Prot
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		{
			name: "Simple test",
			fields: fields{
				prot: []Prot{
					{id: "BLABLA", desc: "Something", seq: "ACDEFGH"},
				},
			},
			wantWriter: ">BLABLA\tSomething\nACDEFGH\n",
			wantErr:    false,
		},
		{
			name: "Full test",
			fields: fields{
				prot: []Prot{
					{id: "BLABLA", desc: "Something", seq: "ACDEFGH"},
					{id: "BLA2", desc: "XXX", seq: "ABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ"},
					{id: "BLA3", desc: "XXX", seq: "ABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ"},
					{id: "BLA4", desc: "ZZZ", seq: "ABCDEFGH"},
				},
			},
			wantWriter: ">BLABLA\tSomething\nACDEFGH\n>BLA2\tXXX\nABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ\nABCDEFGHIJ\n>BLA3\tXXX\nABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ\n>BLA4\tZZZ\nABCDEFGH\n",
			wantErr:    false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fasta{
				prot: tt.fields.prot,
			}
			writer := &bytes.Buffer{}
			if err := f.Write(writer); (err != nil) != tt.wantErr {
				t.Errorf("Fasta.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Fasta.Write() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
