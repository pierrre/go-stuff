// Package passwordentropy calculates the entropy of a password.
package passwordentropy

import (
	"math"
)

// Calculate returns the entropy (in bits) of the given password string.
//
// The entropy is estimated as L * log2(N): L is the password length (in runes)
// and N is the size of the character pool the password is assumed to be drawn
// from.
//
// The pool is the sum of the character classes present in the password.
// Digits (0-9) contribute 10. Letters contribute 26 per case (lower and/or
// upper), unless the only letters used are in the hex range (a-f, A-F), in
// which case they contribute 6 per case: this avoids overestimating the
// entropy of hex strings such as tokens from encoding/hex or
// crypto/rand.Text. Any other character contributes 1 per distinct rune, so a
// password made of a single repeated symbol has an entropy of 0.
//
// An empty password returns 0.
func Calculate(s string) float64 {
	var c classes
	for _, r := range s {
		c.add(r)
	}
	n := c.count()
	if n == 0 {
		return 0
	}
	return float64(c.length) * math.Log2(float64(n))
}

// classes tracks which character classes appear in a password.
type classes struct {
	length    int
	digit     bool
	hexLetter bool // a hex letter (a-f or A-F) is present.
	letter    bool // a non-hex letter (g-z or G-Z) is present.
	lower     bool
	upper     bool
	others    map[rune]struct{}
}

func (c *classes) add(r rune) {
	c.length++
	switch {
	case inRange(r, '0', '9'):
		c.digit = true
	case inRange(r, 'a', 'f'):
		c.hexLetter = true
		c.lower = true
	case inRange(r, 'A', 'F'):
		c.hexLetter = true
		c.upper = true
	case inRange(r, 'a', 'z'):
		c.letter = true
		c.lower = true
	case inRange(r, 'A', 'Z'):
		c.letter = true
		c.upper = true
	default:
		if c.others == nil {
			c.others = make(map[rune]struct{})
		}
		c.others[r] = struct{}{}
	}
}

func (c *classes) count() int {
	n := 0
	if c.digit {
		n += 10
	}
	ln := 0
	if c.letter {
		ln = 26
	} else if c.hexLetter {
		ln = 6
	}
	if c.lower {
		n += ln
	}
	if c.upper {
		n += ln
	}
	n += len(c.others)
	return n
}

func inRange(r, lo, hi rune) bool {
	return r >= lo && r <= hi
}
