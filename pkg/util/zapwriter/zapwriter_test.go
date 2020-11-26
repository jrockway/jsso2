package zapwriter

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestWriter(t *testing.T) {
	buf := new(zaptest.Buffer)
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
	var got []string
	for i := 0; i < 100; i++ {
		got = buf.Lines()
		if len(got) > 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if diff := cmp.Diff(got, []string{`{"msg":"foobar"}`}); diff != "" {
		t.Errorf("output:\n%s", diff)
	}
}
