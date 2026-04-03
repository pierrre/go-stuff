//go:build darwin

package lowpriority

import (
	"golang.org/x/sys/unix"
)

func setThreadLowPriority() error {
	return unix.Setpriority(3, 0, 19) //nolint:wrapcheck // It's OK.
}
