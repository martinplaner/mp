package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"
)

// Generator generates compound words.
type Generator struct {
	vocabulary map[rune][]string
	delimiter  string
}

// Generate generates a compound word for the given term.
func (g *Generator) Generate(term string) (string, error) {
	var words []string
	for _, rune := range term {
		pool := g.vocabulary[rune]
		n := len(pool)
		if n < 1 {
			return "", fmt.Errorf("could not find words for rune %q in vocabulary", rune)
		}
		words = append(words, pool[rand.Intn(n)])
	}
	return strings.Join(words, g.delimiter), nil
}

// GeneratorFromFile creates a new Generator with a vocabulary loaded from the given file path.
func GeneratorFromFile(path, delimiter string) (*Generator, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %v: %w", path, err)
	}
	defer file.Close()

	vocabulary := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if utf8.RuneCountInString(line) < 1 || strings.HasPrefix(line, "#") {
			continue
		}
		rune, _ := utf8.DecodeRuneInString(line)
		if rune == utf8.RuneError {
			return nil, fmt.Errorf("error decoding starting rune in %q", line)
		}
		vocabulary[rune] = append(vocabulary[rune], line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input %v: %w", path, err)
	}

	return &Generator{
		vocabulary: vocabulary,
		delimiter:  delimiter,
	}, nil
}
