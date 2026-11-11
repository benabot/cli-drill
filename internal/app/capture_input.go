package app

import (
	"errors"
	"io"
	"time"

	"github.com/benabot/cli-drill/internal/exercise"
)

const escapeGap = 50 * time.Millisecond

// A zero wait blocks for a key. A positive wait returns EOF on an inter-byte
// timeout. An empty key means the input was consumed but isn't a drill answer.
func readCaptureKey(next func(time.Duration) (byte, error)) (string, error) {
	b, err := next(0)
	if err != nil {
		return "", err
	}
	switch b {
	case '\x1b':
		return readCaptureEscape(next)
	case '\r', '\n':
		return "Enter", nil
	case 0x1f, 0x7f:
		// The corpus accepts these literal bytes. Display names are separate.
		return string(b), nil
	}
	if notation, ok := exercise.KeyByteNotation(b); ok {
		return notation, nil
	}
	if b >= ' ' && b <= '~' {
		return string(b), nil
	}
	return "", nil
}

func readCaptureEscape(next func(time.Duration) (byte, error)) (string, error) {
	b, err := next(escapeGap)
	for err == nil && b == '\x1b' {
		// Alt-modified special keys can have an extra ESC prefix.
		b, err = next(escapeGap)
	}
	if errors.Is(err, io.EOF) {
		return "Esc", nil
	}
	if err != nil {
		return "", err
	}
	switch b {
	case '\x03':
		return "Ctrl+C", nil
	case '[', 'O':
		err = discardCaptureCSI(next, b == '[')
	case ']', 'P', 'X', '^', '_':
		err = discardCaptureString(next, b == ']')
	default:
		// A Meta key, or an escape sequence with intermediate bytes.
		for err == nil && b >= 0x20 && b <= 0x2f {
			b, err = next(escapeGap)
		}
	}
	if errors.Is(err, io.EOF) {
		return "", nil
	}
	return "", err
}

func discardCaptureCSI(next func(time.Duration) (byte, error), csi bool) error {
	const pasteStart = "200~"
	var prefix [len(pasteStart)]byte
	for i := 0; ; i++ {
		b, err := next(escapeGap)
		if err != nil {
			return err
		}
		if i < len(prefix) {
			prefix[i] = b
		}
		if csi && i == 0 && b == '[' {
			// Linux console F1–F5 use ESC [ [ followed by a final byte.
			continue
		}
		if b >= 0x40 && b <= 0x7e {
			if csi && i == len(prefix)-1 && string(prefix[:]) == pasteStart {
				return discardCapturePaste(next)
			}
			return nil
		}
	}
}

func discardCaptureString(next func(time.Duration) (byte, error), osc bool) error {
	escaped := false
	for {
		b, err := next(escapeGap)
		if err != nil {
			return err
		}
		if (osc && b == '\a') || (escaped && b == '\\') {
			return nil
		}
		escaped = b == '\x1b'
	}
}

func discardCapturePaste(next func(time.Duration) (byte, error)) error {
	const end = "\x1b[201~"
	matched := 0
	for {
		b, err := next(escapeGap)
		if err != nil {
			return err
		}
		if b == end[matched] {
			matched++
			if matched == len(end) {
				return nil
			}
		} else if b == end[0] {
			matched = 1
		} else {
			matched = 0
		}
	}
}

func capturedKeyLabel(key string) string {
	switch key {
	case "\x7f":
		return "DEL (Retour arrière)"
	case "\x1f":
		return "Ctrl+_"
	case "Ctrl+I":
		return "Tab (Ctrl+I)"
	default:
		return key
	}
}
