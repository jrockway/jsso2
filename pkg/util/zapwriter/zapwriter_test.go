package zapwriter

import (
	"bytes"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type buffer struct {
	c chan []byte
}

func (w *buffer) Write(b []byte) (int, error) {
	w.c <- b
	return len(b), nil
}

func (w *buffer) Sync() error {
	return nil
}

func TestWriter(t *testing.T) {
	buf := &buffer{c: make(chan []byte)}
	c := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
		}),
		buf,
		zapcore.DebugLevel,
	)
	l := zap.New(c)
	w := New(l)
	if n, err := w.Write([]byte("foo")); err != nil {
		t.Fatalf("foo: %v", err)
	} else if n != 3 {
		t.Fatalf("foo: short write %d", n)
	}
	if n, err := w.Write([]byte("bar\n")); err != nil {
		t.Fatalf("bar: %v", err)
	} else if n != 4 {
		t.Fatalf("bar: short write %d", n)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	var got []byte
	select {
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for write")
	case got = <-buf.c:
	}

	if want := []byte(`{"msg":"foobar"}` + "\n"); !bytes.Equal(got, want) {
		t.Errorf("line:\n  got: %q\n want: %q", got, want)
	}
}
