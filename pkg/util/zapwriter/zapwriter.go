// Package zapwriter implements an io.Writer that logs one line at a time to zap.
package zapwriter

import (
	"bufio"
	"io"

	"go.uber.org/zap"
)

func New(l *zap.Logger) io.WriteCloser {
	r, w := io.Pipe()
	s := bufio.NewScanner(r)
	go func() {
		for s.Scan() {
			l.Debug(s.Text())
		}
		if err := s.Err(); err != nil {
			l.Error("zapwriter exiting early", zap.Error(err))
			r.CloseWithError(err)
			return
		}
		r.Close()
	}()
	return w
}
