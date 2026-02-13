// Package passwordentropy calculates the entropy of a password.
package passwordentropy

import (
	"math"
)

// Calculate returns the entropy of the given password string.
// The entropy is calculated based on the length of the password and the number of possible characters used in the password.
//
//nolint:gocyclo // The function is simple enough and the complexity is not a problem.
func Calculate(s string) float64 {
	hasDigit := false
	hasHexLetter := false
	hasLetter := false
	hasLower := false
	hasUpper := false
	var others map[rune]struct{}
	l := 0
	for _, c := range s {
		l++
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case c >= 'a' && c <= 'f':
			hasHexLetter = true
			hasLower = true
		case c >= 'A' && c <= 'F':
			hasHexLetter = true
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLetter = true
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasLetter = true
			hasUpper = true
		default:
			if others == nil {
				others = make(map[rune]struct{})
			}
			others[c] = struct{}{}
		}
	}
	n := 0
	if hasDigit {
		n += 10
	}
	ln := 0
	if hasLetter {
		ln = 26
	} else if hasHexLetter {
		ln = 6
	}
	if hasLower {
		n += ln
	}
	if hasUpper {
		n += ln
	}
	n += len(others)
	if n > 0 {
		return float64(l) * math.Log2(float64(n))
	}
	return 0
}
