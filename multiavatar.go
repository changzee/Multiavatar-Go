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
	// selectedTheme forces theme letter ("A","B","C") for all parts if set
	selectedTheme *string
	// forcePartV allows overriding the part version (e.g., "eyes":"07")
	forcePartV map[string]string
	// allowedVersions restricts each part to a set of allowed versions; selection is deterministic within the set
	allowedVersions map[string][]string
	// per-part theme override: e.g., {"eyes":"B"}
	partTheme map[string]string
	// allowed theme letters per part: e.g., {"top": {"A","C"}}
	allowedThemes map[string][]string
	// disable specific parts: {"top": true} to skip rendering that part
	disabledParts map[string]bool
	// overrideColors allows overriding the colors array for a specific part
	// e.g., {"head": {"#f2c280"}} to force skin tone
	overrideColors map[string][]string
}

// Option is a function that configures a generation option.
type Option func(*config)

// WithTheme forces the theme letter ("A","B","C") globally.
func WithTheme(theme string) Option {
	return func(c *config) {
		t := strings.ToUpper(strings.TrimSpace(theme))
		if t == "A" || t == "B" || t == "C" {
			c.selectedTheme = &t
		}
	}
}

// WithPartVersion forces a specific part to use a given version "00".."15".
func WithPartVersion(partName, partVersion string) Option {
	return func(c *config) {
		if c.forcePartV == nil {
			c.forcePartV = make(map[string]string)
		}
		pn := strings.TrimSpace(partName)
		pv := strings.TrimSpace(partVersion)
		// basic validation: partName must be one of known parts and version must be 2-digit
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			if len(pv) == 2 {
				c.forcePartV[pn] = pv
			}
		}
	}
}

// WithPartColors overrides the colors array used for a specific part.
// For example, WithPartColors("head", []string{"#f2c280"}) to set skin tone.
func WithPartColors(partName string, colors []string) Option {
	return func(c *config) {
		if c.overrideColors == nil {
			c.overrideColors = make(map[string][]string)
		}
		pn := strings.TrimSpace(partName)
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			// store a copy to avoid external mutation
			cp := make([]string, len(colors))
			for i := range colors {
				cp[i] = strings.TrimSpace(colors[i])
			}
			c.overrideColors[pn] = cp
		}
	}
}

// Convenience options for common cases

// WithSkinColor sets the head (skin) primary color.
func WithSkinColor(hex string) Option {
	return WithPartColors("head", []string{strings.TrimSpace(hex)})
}

// WithEyesColors sets the eyes colors array (primary, secondary, etc.).
func WithEyesColors(colors ...string) Option {
	return WithPartColors("eyes", colors)
}

// WithTopColors sets the hair/top colors array.
func WithTopColors(colors ...string) Option {
	return WithPartColors("top", colors)
}

// WithEnvColor sets the environment/background circle color.
func WithEnvColor(hex string) Option {
	return WithPartColors("env", []string{strings.TrimSpace(hex)})
}

// WithClothesColors sets clothes colors array.
func WithClothesColors(colors ...string) Option {
	return WithPartColors("clo", colors)
}

// WithMouthColors sets mouth colors array.
func WithMouthColors(colors ...string) Option {
	return WithPartColors("mouth", colors)
}

// WithPartTheme forces theme letter ("A","B","C") for a specific part.
func WithPartTheme(partName, theme string) Option {
	return func(c *config) {
		if c.partTheme == nil {
			c.partTheme = make(map[string]string)
		}
		pn := strings.TrimSpace(partName)
		t := strings.ToUpper(strings.TrimSpace(theme))
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			if t == "A" || t == "B" || t == "C" {
				c.partTheme[pn] = t
			}
		}
	}
}

// WithAllowedThemes restricts a part to given theme letters (e.g., ["A","C"]).
// Deterministic selection within the list based on the input hash slice.
func WithAllowedThemes(partName string, themesList []string) Option {
	return func(c *config) {
		if c.allowedThemes == nil {
			c.allowedThemes = make(map[string][]string)
		}
		pn := strings.TrimSpace(partName)
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			var tl []string
			for _, t := range themesList {
				tu := strings.ToUpper(strings.TrimSpace(t))
				if tu == "A" || tu == "B" || tu == "C" {
					tl = append(tl, tu)
				}
			}
			if len(tl) > 0 {
				c.allowedThemes[pn] = tl
			}
		}
	}
}

// WithoutPart disables rendering a specific part (e.g., "top" to remove hair).
func WithoutPart(partName string) Option {
	return func(c *config) {
		if c.disabledParts == nil {
			c.disabledParts = make(map[string]bool)
		}
		pn := strings.TrimSpace(partName)
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			c.disabledParts[pn] = true
		}
	}
}

// WithAllowedVersions restricts a part to given version list (e.g., ["01","03","07"]).
// The algorithm will pick deterministically within the list based on the input hash.
func WithAllowedVersions(partName string, versions []string) Option {
	return func(c *config) {
		if c.allowedVersions == nil {
			c.allowedVersions = make(map[string][]string)
		}
		pn := strings.TrimSpace(partName)
		switch pn {
		case "env", "clo", "head", "mouth", "eyes", "top":
			// sanitize to 2-digit codes
			var vlist []string
			for _, v := range versions {
				v = strings.TrimSpace(v)
				if len(v) == 2 {
					vlist = append(vlist, v)
				}
			}
			if len(vlist) > 0 {
				c.allowedVersions[pn] = vlist
			}
		}
	}
}

// Convenience: restrict common parts
func WithAllowedHeadVersions(versions ...string) Option { return WithAllowedVersions("head", versions) }
func WithAllowedEyesVersions(versions ...string) Option { return WithAllowedVersions("eyes", versions) }
func WithAllowedTopVersions(versions ...string) Option  { return WithAllowedVersions("top", versions) }

// WithGender applies preset style filters for male/female/unisex.
// It restricts allowed versions/themes for certain parts to achieve gendered styling
// while keeping deterministic selection within those sets.
func WithGender(gender string) Option {
	return func(c *config) {
		g := strings.ToLower(strings.TrimSpace(gender))

		// ensure maps
		if c.allowedVersions == nil {
			c.allowedVersions = make(map[string][]string)
		}
		if c.allowedThemes == nil {
			c.allowedThemes = make(map[string][]string)
		}
		if c.partTheme == nil {
			c.partTheme = make(map[string]string)
		}

		switch g {
		case "female", "woman", "girl", "f", "♀":
			// Presets偏向女性风格
			c.allowedVersions["top"] = []string{"01", "03", "07", "10"}
			c.allowedVersions["eyes"] = []string{"03", "11"}
			c.allowedThemes["top"] = []string{"A", "C"}
			// 加强女性风格的主题倾向
			c.partTheme["top"] = "C"
			c.partTheme["eyes"] = "C"
		case "male", "man", "boy", "m", "♂":
			// Presets偏向男性风格
			c.allowedVersions["top"] = []string{"04", "05", "14"}
			c.allowedVersions["eyes"] = []string{"09", "10"}
			c.allowedThemes["top"] = []string{"A", "B"}
			// 加强男性风格的主题倾向
			c.partTheme["top"] = "B"
			c.partTheme["eyes"] = "B"
		default:
			// Unisex更广的集合
			c.allowedVersions["top"] = []string{"01", "03", "04", "05", "07", "10", "14"}
			c.allowedVersions["eyes"] = []string{"03", "09", "10", "11"}
			c.allowedThemes["top"] = []string{"A", "B", "C"}
			// 中性不强制具体主题，让 allowedThemes 生效
		}
	}
}

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

	// ensure internal maps are initialized
	if cfg.forcePartV == nil {
		cfg.forcePartV = make(map[string]string)
	}
	if cfg.allowedVersions == nil {
		cfg.allowedVersions = make(map[string][]string)
	}
	if cfg.partTheme == nil {
		cfg.partTheme = make(map[string]string)
	}
	if cfg.allowedThemes == nil {
		cfg.allowedThemes = make(map[string][]string)
	}
	if cfg.disabledParts == nil {
		cfg.disabledParts = make(map[string]bool)
	}
	if cfg.overrideColors == nil {
		cfg.overrideColors = make(map[string][]string)
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

		// Apply forced/global/per-part theme/version if configured
		if cfg.selectedTheme != nil {
			theme = *cfg.selectedTheme
		}
		if pt, ok := cfg.partTheme[name]; ok {
			theme = pt
		} else if allowedT, ok := cfg.allowedThemes[name]; ok && len(allowedT) > 0 {
			theme = allowedT[val%len(allowedT)]
		}

		if forced, ok := cfg.forcePartV[name]; ok && len(forced) == 2 {
			partV = forced
		} else if allowed, ok := cfg.allowedVersions[name]; ok && len(allowed) > 0 {
			partV = allowed[val%len(allowed)]
		}

		// 4d. Get the final SVG part with colors, allowing overrides
		selectedParts[name] = getFinalPartWithOverride(name, partV, theme, cfg.overrideColors[name])
	}

	// 5. Assemble the final SVG
	var finalSVG strings.Builder
	finalSVG.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 231 231">`)

	if !cfg.withoutBackground && !cfg.disabledParts["env"] {
		finalSVG.WriteString(selectedParts["env"])
	}
	if !cfg.disabledParts["head"] {
		finalSVG.WriteString(selectedParts["head"])
	}
	if !cfg.disabledParts["clo"] {
		finalSVG.WriteString(selectedParts["clo"])
	}
	if !cfg.disabledParts["top"] {
		finalSVG.WriteString(selectedParts["top"])
	}
	if !cfg.disabledParts["eyes"] {
		finalSVG.WriteString(selectedParts["eyes"])
	}
	if !cfg.disabledParts["mouth"] {
		finalSVG.WriteString(selectedParts["mouth"])
	}

	finalSVG.WriteString(`</svg>`)

	return finalSVG.String()
}

// getFinalPartWithOverride retrieves the raw SVG string for a part,
// and replaces color placeholders, allowing optional color overrides.
func getFinalPartWithOverride(partName, partV, theme string, override []string) string {
	colors, ok := themes[partV][theme][partName]
	if !ok {
		return "" // Should not happen with correct logic
	}

	// If override provided, use it (truncate/extend matching placeholders count on use)
	if override != nil && len(override) > 0 {
		// copy to avoid mutating themes
		cp := make([]string, len(override))
		copy(cp, override)
		colors = cp
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
