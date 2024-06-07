package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/mozillazg/go-unidecode"

	"strings"
)

// NormalizeText normalizes the text by converting it to lowercase and removing punctuation
func NormalizeText(text string) string {
	text = strings.ToLower(text)
	text = unidecode.Unidecode(text)
	replacer := strings.NewReplacer(".", "", ",", "", "!", "", "?", "", ":", "", ";", "", "'", "", "\"", "")
	return replacer.Replace(text)
}

// GenerateHash generates a SHA256 hash of the normalized text
func GenerateHash(text string) string {
	normalizedText := NormalizeText(text)
	hash := sha256.Sum256([]byte(normalizedText))
	return hex.EncodeToString(hash[:])
}

func CalculateCosineSimilarity(s1, s2 string) float64 {
	m := metrics.NewJaroWinkler()
	return strutil.Similarity(s1, s2, m)

}
