package multiavatar

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// config holds the configuration for generating an avatar.
type config struct {
	withoutBackground bool
}

// Option is a function that configures a generation option.
type Option func(*config)

// WithoutBackground is an option to generate an avatar with a transparent background.
func WithoutBackground() Option {
	return func(c *config) {
		c.withoutBackground = true
	}
}

// Generate creates an SVG avatar string from an input string based on a deterministic algorithm.
// It is thread-safe.
func Generate(input string, opts ...Option) string {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	if input == "" {
		return ""
	}

	// 1. SHA-256 hash
	hashBytes := sha256.Sum256([]byte(input))
	hexHash := hex.EncodeToString(hashBytes[:])

	// 2. Remove non-digits (mimicking JS replace(/\D/g, ''))
	re := regexp.MustCompile(`[^0-9]`)
	sha256Numbers := re.ReplaceAllString(hexHash, "")

	// 3. Get the first 12 digits
	hashStr := sha256Numbers
	if len(hashStr) > 12 {
		hashStr = hashStr[:12]
	}

	// 4. Determine parts
	partNames := []string{"env", "clo", "head", "mouth", "eyes", "top"}
	selectedParts := make(map[string]string)

	for i, name := range partNames {
		// 4a. Take 2 digits
		valStr := hashStr[i*2 : i*2+2]
		val, _ := strconv.Atoi(valStr)

		// 4b. Scale to 0-47 range
		nr := int(math.Round(float64(val) * 47 / 100))

		// 4c. Determine version (partV) and theme (A, B, C)
		var partV, theme string
		if nr > 31 {
			partV = fmt.Sprintf("%02d", nr-32)
			theme = "C"
		} else if nr > 15 {
			partV = fmt.Sprintf("%02d", nr-16)
			theme = "B"
		} else {
			partV = fmt.Sprintf("%02d", nr)
			theme = "A"
		}

		// 4d. Get the final SVG part with colors
		selectedParts[name] = getFinalPart(name, partV, theme)
	}

	// 5. Assemble the final SVG
	var finalSVG strings.Builder
	finalSVG.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 231 231">`)

	if !cfg.withoutBackground {
		finalSVG.WriteString(selectedParts["env"])
	}
	finalSVG.WriteString(selectedParts["head"])
	finalSVG.WriteString(selectedParts["clo"])
	finalSVG.WriteString(selectedParts["top"])
	finalSVG.WriteString(selectedParts["eyes"])
	finalSVG.WriteString(selectedParts["mouth"])

	finalSVG.WriteString(`</svg>`)

	return finalSVG.String()
}

// getFinalPart retrieves the raw SVG string for a part, and replaces color placeholders.
func getFinalPart(partName, partV, theme string) string {
	colors, ok := themes[partV][theme][partName]
	if !ok {
		return "" // Should not happen with correct logic
	}

	partID, _ := strconv.Atoi(partV)
	partIndex := map[string]int{"env": 0, "clo": 1, "head": 2, "mouth": 3, "eyes": 4, "top": 5}[partName]

	svgString := parts[partID][partIndex]

	// Replace color placeholders like "#01;"
	re := regexp.MustCompile(`#(.*?);`)
	matches := re.FindAllString(svgString, -1)

	resultFinal := svgString
	if matches != nil {
		for i, placeholder := range matches {
			if i < len(colors) {
				resultFinal = strings.Replace(resultFinal, placeholder, colors[i]+";", 1)
			}
		}
	}

	return resultFinal
}
