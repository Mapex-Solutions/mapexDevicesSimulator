package services

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// placeholderRe matches a {{ ... }} token.
var placeholderRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Render resolves the {{...}} placeholders in a template at send time. Unknown
// tokens are left untouched so a typo is visible rather than silently dropped.
func Render(template, deviceID string, counter int64) string {
	return placeholderRe.ReplaceAllStringFunc(template, func(match string) string {
		token := strings.TrimSpace(match[2 : len(match)-2])
		return resolvePlaceholder(token, deviceID, counter)
	})
}

// resolvePlaceholder maps a single token to its value.
func resolvePlaceholder(token, deviceID string, counter int64) string {
	switch {
	case token == "now":
		return time.Now().UTC().Format(time.RFC3339)
	case token == "counter":
		return strconv.FormatInt(counter, 10)
	case token == "deviceId":
		return deviceID
	case token == "uuid":
		return uuid.NewString()
	case strings.HasPrefix(token, "randInt(") && strings.HasSuffix(token, ")"):
		return renderRandInt(token[len("randInt(") : len(token)-1])
	case strings.HasPrefix(token, "randFloat(") && strings.HasSuffix(token, ")"):
		return renderRandFloat(token[len("randFloat(") : len(token)-1])
	default:
		return "{{" + token + "}}"
	}
}

// renderRandInt resolves randInt(min,max) to an integer in [min,max].
func renderRandInt(args string) string {
	lo, hi, ok := parseRange(args)
	if !ok {
		return "0"
	}
	min, max := int(lo), int(hi)
	if max < min {
		min, max = max, min
	}
	return strconv.Itoa(min + rand.Intn(max-min+1))
}

// renderRandFloat resolves randFloat(min,max) to a decimal in [min,max].
func renderRandFloat(args string) string {
	min, max, ok := parseRange(args)
	if !ok {
		return "0"
	}
	if max < min {
		min, max = max, min
	}
	return strconv.FormatFloat(min+rand.Float64()*(max-min), 'f', 2, 64)
}

// parseRange parses a "min,max" argument list into two floats.
func parseRange(args string) (float64, float64, bool) {
	parts := strings.Split(args, ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	a, errA := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	b, errB := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if errA != nil || errB != nil {
		return 0, 0, false
	}
	return a, b, true
}
