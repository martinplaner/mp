package main

import (
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Generator
		term    string
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			g:       &Generator{},
			term:    "",
			want:    "",
			wantErr: false,
		},
		{
			name: "Alpha-Bravo",
			g: &Generator{
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
			g:       &Generator{},
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
		args    args
		wantErr bool
	}{
		{
			name: "full word list",
			args: args{
				path:      "words.txt",
				delimiter: "-",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratorFromFile(tt.args.path, tt.args.delimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratorFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			s, err := got.Generate("ABC")
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
