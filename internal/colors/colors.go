package colors

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Palette holds all generated CSS color variables.
type Palette struct {
	Primary           string
	PrimaryHover      string
	PrimaryLight      string
	BgMain            string
	BgCard            string
	TextMain          string
	TextSecondary     string
	TextMuted         string
	Border            string
	Success           string
	Error             string
	Warning           string
	DarkBgMain        string
	DarkBgCard        string
	DarkBgInput       string
	DarkTextMain      string
	DarkTextSecondary string
	DarkBorder        string
	Hue               int
}

// Generate creates a cohesive HSL palette from a random hue.
// Skips red-ish range for a professional look (matches bash colors.sh).
func Generate() Palette {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	hue := rng.Intn(300) + 20

	return Palette{
		Primary:           hslToHex(float64(hue), 70, 55),
		PrimaryHover:      hslToHex(float64(hue), 70, 45),
		PrimaryLight:      hslToHex(float64(hue), 70, 70),
		BgMain:            hslToHex(float64(hue), 20, 96),
		BgCard:            hslToHex(float64(hue), 15, 99),
		TextMain:          hslToHex(float64(hue), 20, 10),
		TextSecondary:     hslToHex(float64(hue), 15, 25),
		TextMuted:         hslToHex(float64(hue), 10, 45),
		Border:            hslToHex(float64(hue), 20, 82),
		Success:           hslToHex(140, 60, 45),
		Error:             hslToHex(0, 70, 55),
		Warning:           hslToHex(38, 90, 55),
		DarkBgMain:        hslToHex(float64(hue), 30, 8),
		DarkBgCard:        hslToHex(float64(hue), 25, 13),
		DarkBgInput:       hslToHex(float64(hue), 20, 20),
		DarkTextMain:      hslToHex(float64(hue), 15, 96),
		DarkTextSecondary: hslToHex(float64(hue), 12, 78),
		DarkBorder:        hslToHex(float64(hue), 18, 28),
		Hue:               hue,
	}
}

// InjectIntoCSS applies the palette variables into a style.css file.
func (p Palette) InjectIntoCSS(cssPath string) error {
	data, err := os.ReadFile(cssPath)
	if err != nil {
		return err
	}

	replacements := [][2]string{
		{"--primary:", fmt.Sprintf("--primary:           %s;", p.Primary)},
		{"--primary-hover:", fmt.Sprintf("--primary-hover:     %s;", p.PrimaryHover)},
		{"--primary-light:", fmt.Sprintf("--primary-light:     %s;", p.PrimaryLight)},
		{"--bg-main:", fmt.Sprintf("--bg-main:           %s;", p.BgMain)},
		{"--bg-card:", fmt.Sprintf("--bg-card:           %s;", p.BgCard)},
		{"--text-main:", fmt.Sprintf("--text-main:         %s;", p.TextMain)},
		{"--text-secondary:", fmt.Sprintf("--text-secondary:    %s;", p.TextSecondary)},
		{"--text-muted:", fmt.Sprintf("--text-muted:        %s;", p.TextMuted)},
		{"--border:", fmt.Sprintf("--border:            %s;", p.Border)},
		{"--success:", fmt.Sprintf("--success:           %s;", p.Success)},
		{"--error:", fmt.Sprintf("--error:             %s;", p.Error)},
		{"--warning:", fmt.Sprintf("--warning:           %s;", p.Warning)},
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, pair := range replacements {
			if strings.HasPrefix(trimmed, pair[0]) {
				indent := lineIndent(line)
				lines[i] = indent + pair[1]
				break
			}
		}
	}

	return os.WriteFile(cssPath, []byte(strings.Join(lines, "\n")), 0644)
}

// ── Internal math ─────────────────────────────────────────────────────────────

func hslToHex(h, s, l float64) string {
	h /= 360
	s /= 100
	l /= 100

	if s == 0 {
		v := int(l * 255)
		return fmt.Sprintf("#%02x%02x%02x", v, v, v)
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p2 := 2*l - q

	r := hue2rgb(p2, q, h+1.0/3)
	g := hue2rgb(p2, q, h)
	b := hue2rgb(p2, q, h-1.0/3)

	return fmt.Sprintf("#%02x%02x%02x", clampByte(r), clampByte(g), clampByte(b))
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

func clampByte(v float64) int {
	n := int(v * 255)
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return n
}

func lineIndent(s string) string {
	for i, c := range s {
		if c != ' ' && c != '\t' {
			return s[:i]
		}
	}
	return ""
}
