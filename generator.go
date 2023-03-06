package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Generator generates words from a given term.
type Generator interface {
	Generate(term string) (string, error)
}

// CompoundGenerator generates compound words.
type CompoundGenerator struct {
	vocabulary map[rune][]string
	delimiter  string
}

// Generate generates a compound word for the given term.
func (g *CompoundGenerator) Generate(term string) (string, error) {
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

// CompoundGeneratorFromFile creates a new compound word Generator with a vocabulary loaded from the given file path.
func CompoundGeneratorFromFile(path, delimiter string) (Generator, error) {
	lines, err := loadLinesFromFile(path)

	if err != nil {
		return nil, fmt.Errorf("loading lines from file %v: %w", path, err)
	}

	vocabulary := make(map[rune][]string)
	for _, line := range lines {
		rune, _ := utf8.DecodeRuneInString(line)
		if rune == utf8.RuneError {
			return nil, fmt.Errorf("error decoding starting rune in %q", line)
		}
		// normalize runes to upper case
		rune = unicode.ToUpper(rune)
		vocabulary[rune] = append(vocabulary[rune], line)
	}

	return &CompoundGenerator{
		vocabulary: vocabulary,
		delimiter:  delimiter,
	}, nil
}

// AdjectiveGenerator generates a list of adjectives and a final noun.
type AdjectiveGenerator struct {
	adjectives map[rune][]string
	nouns      map[rune][]string
	delimiter  string
}

// Generate generates a compound word for the given term.
func (g *AdjectiveGenerator) Generate(term string) (string, error) {
	var words []string
	for i, rune := range term {
		pool := g.adjectives[rune]

		// use a noun for the last rune of term
		if i >= len(term)-1 {
			pool = g.nouns[rune]
		}

		n := len(pool)
		if n < 1 {
			return "", fmt.Errorf("could not find words for rune %q in vocabulary", rune)
		}
		words = append(words, pool[rand.Intn(n)])
	}
	return strings.Join(words, g.delimiter), nil
}

// AdjectiveGeneratorFromFile creates a new adjective word Generator with a vocabulary loaded from the given file path.
func AdjectiveGeneratorFromFile(path, delimiter string) (Generator, error) {
	lines, err := loadLinesFromFile(path)

	if err != nil {
		return nil, fmt.Errorf("loading lines from file %v: %w", path, err)
	}

	adjectives := make(map[rune][]string)
	nouns := make(map[rune][]string)

	for _, line := range lines {
		rune, _ := utf8.DecodeRuneInString(line)
		if rune == utf8.RuneError {
			return nil, fmt.Errorf("error decoding starting rune in %q", line)
		}
		// normalize runes to upper case
		rune = unicode.ToUpper(rune)
		split := strings.Split(line, ",")

		if len(split) == 1 {
			nouns[rune] = append(nouns[rune], split[0])
		} else if len(split) > 1 && split[1] == "a" {
			adjectives[rune] = append(adjectives[rune], split[0])
		} else {
			return nil, fmt.Errorf("unknown word type: %v", split[1])
		}
	}

	return &AdjectiveGenerator{
		adjectives: adjectives,
		nouns:      nouns,
		delimiter:  delimiter,
	}, nil
}

func loadLinesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %v: %w", path, err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if utf8.RuneCountInString(line) < 1 || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input %v: %w", path, err)
	}

	return lines, nil
}
