package passwordentropy

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

func Test(t *testing.T) {
	for _, s := range []string{
		"",
		"abc",
		"ABC",
		"Abc",
		"a1b2c3",
		"A1B2C3",
		"123456",
		"password",
		"Password",
		"Password1",
		"Password1!",
		",?;.:/!&",
		"3TIFXL2JL3QVDWC6JBB2YO6GWI", // generated from [crypto/rand.Text]
	} {
		e := Calculate(s)
		assertauto.Equal(t, s)
		assertauto.Equal(t, e)
	}
}
