//go:build unix

package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/unix"
)

const rawCaptureSupported = true

func readRawKey(file *os.File) (_ string, result error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTSTP)
	defer func() {
		signal.Stop(signals)
		select {
		case sig := <-signals:
			// A signal can arrive after the read, while raw mode is restored.
			result = errors.Join(result, captureSignalError{sig})
		default:
		}
	}()

	fd := file.Fd()
	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("activate raw keyboard capture: %w", err)
	}
	defer func() {
		if err := term.Restore(fd, state); err != nil {
			result = errors.Join(result, fmt.Errorf("restore terminal after keyboard capture: %w", err))
		}
	}()

	fds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}
	return readCaptureKey(func(wait time.Duration) (byte, error) {
		deadline := time.Now().Add(wait)
		for {
			select {
			case sig := <-signals:
				return 0, captureSignalError{sig}
			default:
			}
			// Bound signal latency without a background reader that could
			// consume input after switching back to canonical mode.
			const signalInterval = 50 * time.Millisecond
			pollFor := signalInterval
			if wait > 0 {
				remaining := time.Until(deadline)
				if remaining <= 0 {
					return 0, io.EOF
				}
				pollFor = min(pollFor, remaining)
			}
			ready, err := unix.Poll(fds, int((pollFor+time.Millisecond-1)/time.Millisecond))
			if errors.Is(err, unix.EINTR) {
				continue
			}
			if err != nil {
				return 0, err
			}
			if ready == 0 {
				continue
			}
			select {
			case sig := <-signals:
				return 0, captureSignalError{sig}
			default:
			}
			var b [1]byte
			_, err = io.ReadFull(file, b[:])
			return b[0], err
		}
	})
}

type captureSignalError struct{ signal os.Signal }

func (e captureSignalError) Error() string {
	return fmt.Sprintf("capture clavier interrompue (%s)", e.signal)
}

func (e captureSignalError) ExitCode() int {
	// Unix shells expose a signal termination as 128 plus its number.
	return 128 + int(e.signal.(syscall.Signal))
}
