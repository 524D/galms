package elements

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Element
		wantErr bool
	}{
		{
			name: "Simple test elements",
			args: args{r: strings.NewReader(`
[{"Symbol":"Na","Name":"sodium","Number":11,"Isotope":[{"Mass":22.98977,"Abundance":1}]}]
`)},
			want: []Element{
				{
					Symbol: "Na",
					Name:   "sodium",
					Number: 11,
					Isotope: []Isotope{
						{
							Mass: 22.989770, Abundance: 1.0,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.r)
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

func TestInitDefault(t *testing.T) {
	got := InitDefault()
	if got.Elements[0].Name != `hydrogen` {
		t.Errorf("InitDefault() = %s, want %s", got.Elements[0].Name, `hydrogen`)
	}
}
