package utils

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func removeAcentos(input string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, input)
	return result
}

// SanitizeFilename remove acentos, símbolos e normaliza o nome do arquivo
func SanitizeFilename(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	base = removeAcentos(base)

	// Remove caracteres inválidos para nomes de arquivos
	invalid := regexp.MustCompile(`[^\w\s]`)
	base = invalid.ReplaceAllString(base, "")
	base = strings.Join(strings.Fields(base), " ")

	return strings.TrimSpace(base) + ext
}
