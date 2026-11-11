//go:build !unix

package app

import "os"

const rawCaptureSupported = false

func readRawKey(_ *os.File) (string, error) {
	return "", errCaptureUnavailable
}
