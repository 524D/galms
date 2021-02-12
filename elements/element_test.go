package elements

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

// func TestRead1(t *testing.T) {
// 	type args struct {
// 		r io.Reader
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []Element
// 		wantErr bool
// 	}{
// 		{
// 			name: "Simple test elements",
// 			args: args{r: strings.NewReader(`
// [{"Symbol":"Na","Name":"sodium","Number":11,"Isotope":[{"Mass":22.98977,"Abundance":1}]}]
// `)},
// 			want: []Element{
// 				{
// 					Symbol: "Na",
// 					Name:   "sodium",
// 					Number: 11,
// 					Isotope: []Isotope{
// 						{
// 							Mass: 22.989770, Abundance: 1.0,
// 						},
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid JSON",
// 			args: args{r: strings.NewReader(`
// [{"Symbol":Na",'Name':"sodium","Number":11,"Isotope":[{"Mass":22.98977,"Abundance":1}]}]
// `)},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := Read(tt.args.r)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Read() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestInitDefault(t *testing.T) {
	got := InitDefault()
	if got.Elements[0].Name != `hydrogen` {
		t.Errorf("InitDefault() = %s, want %s", got.Elements[0].Name, `hydrogen`)
	}
}

func TestRead(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Elems
		wantErr bool
	}{
		{
			name: "Simple test elements",
			args: args{r: strings.NewReader(`
[{"Symbol":"Na","Name":"sodium","Number":11,"Isotope":[{"Mass":22.98977,"Abundance":1}]}]
`)},
			want: Elems{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			wantErr: false,
		},
		{
			name: "Invalid JSON",
			args: args{r: strings.NewReader(`
[{"Symbol":Na",'Name':"sodium","Number":11,"Isotope":[{"Mass":22.98977,"Abundance":1}]}]
`)},
			want: Elems{
				nil,
				nil,
			},
			wantErr: true,
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

func TestElems_ElemIdx(t *testing.T) {
	type fields struct {
		Elements  []Element
		SymbolMap map[string]int
	}
	type args struct {
		shortName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Simple test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				shortName: "Na",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Fail test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				shortName: "C",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elems{
				Elements:  tt.fields.Elements,
				SymbolMap: tt.fields.SymbolMap,
			}
			got, err := e.ElemIdx(tt.args.shortName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Elems.ElemIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Elems.ElemIdx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElems_Name(t *testing.T) {
	type fields struct {
		Elements  []Element
		SymbolMap map[string]int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 0,
			},
			want:    "sodium",
			wantErr: false,
		},
		{
			name: "Fail test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 666,
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elems{
				Elements:  tt.fields.Elements,
				SymbolMap: tt.fields.SymbolMap,
			}
			got, err := e.Name(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Elems.Name() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Elems.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElems_Symbol(t *testing.T) {
	type fields struct {
		Elements  []Element
		SymbolMap map[string]int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 0,
			},
			want:    "Na",
			wantErr: false,
		},
		{
			name: "Fail test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 666,
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elems{
				Elements:  tt.fields.Elements,
				SymbolMap: tt.fields.SymbolMap,
			}
			got, err := e.Symbol(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Elems.Symbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Elems.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElems_Isotopes(t *testing.T) {
	type fields struct {
		Elements  []Element
		SymbolMap map[string]int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Isotope
		wantErr bool
	}{
		{
			name: "Simple test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 0,
			},
			want: []Isotope{
				{
					Mass: 22.989770, Abundance: 1.0,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail test",
			fields: fields{
				[]Element{
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
				map[string]int{"Na": 0},
			},
			args: args{
				i: 666,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elems{
				Elements:  tt.fields.Elements,
				SymbolMap: tt.fields.SymbolMap,
			}
			got, err := e.Isotopes(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Elems.Isotopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Elems.Isotopes() = %v, want %v", got, tt.want)
			}
		})
	}
}
