package main

import (
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		g       Generator
		term    string
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			g:       &CompoundGenerator{},
			term:    "",
			want:    "",
			wantErr: false,
		},
		{
			name: "Alpha-Bravo",
			g: &CompoundGenerator{
				vocabulary: map[rune][]string{
					rune('A'): {"Alpha"},
					rune('B'): {"Bravo"},
				},
				delimiter: "-",
			},
			term:    "AB",
			want:    "Alpha-Bravo",
			wantErr: false,
		},
		{
			name:    "unknown rune",
			g:       &CompoundGenerator{},
			term:    "A",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			g:       &AdjectiveGenerator{},
			term:    "",
			want:    "",
			wantErr: false,
		},
		{
			name: "alpha-y bravo",
			g: &AdjectiveGenerator{
				adjectives: map[rune][]string{
					rune('A'): {"alpha"},
				},
				nouns: map[rune][]string{
					rune('B'): {"bravo"},
				},
				delimiter: " ",
			},
			term:    "AB",
			want:    "alpha bravo",
			wantErr: false,
		},
		{
			name:    "unknown rune",
			g:       &AdjectiveGenerator{},
			term:    "A",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Generate(tt.term)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Generator.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratorFromFile(t *testing.T) {
	type args struct {
		path      string
		delimiter string
	}
	tests := []struct {
		name    string
		factory func(string, string) (Generator, error)
		args    args
		term    string
		wantErr bool
	}{
		{
			name:    "full word list",
			factory: AdjectiveGeneratorFromFile,
			args: args{
				path:      "words_en.txt",
				delimiter: " ",
			},
			term:    "LB",
			wantErr: false,
		},
		{
			name:    "full word list",
			factory: CompoundGeneratorFromFile,
			args: args{
				path:      "words_de.txt",
				delimiter: "-",
			},
			term:    "MP",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.factory(tt.args.path, tt.args.delimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratorFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			s, err := got.Generate(tt.term)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(s) < 1 {
				t.Errorf("Generate() = %v, want non-empty string", got)
			}
		})
	}
}
