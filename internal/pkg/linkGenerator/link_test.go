package linkGenerator

import (
	"testing"
)

func containsChar(alphabet string, char rune) bool {
	for _, a := range alphabet {
		if a == char {
			return true
		}
	}
	return false
}

func TestGenLink(t *testing.T) {
	tests := []struct {
		name     string
		longLink string
		alphabet string
		length   int
	}{
		{"ShortLinkTest", "www.ozon.ru", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 63},
		{"DifferentAlphabetTest", "e.mail.ru", "abcde", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkHash := NewLinkHash(tt.alphabet, tt.length)
			result := linkHash.GenLink(tt.longLink)

			for _, char := range result {
				if !containsChar(tt.alphabet, char) {
					t.Errorf("Character %c not found in the specified alphabet", char)
				}
			}
		})
	}
}
