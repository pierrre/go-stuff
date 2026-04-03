//go:build !unix && !windows

package lowpriority

func setThreadLowPriority() error {
	return nil // Not implemented.
}
