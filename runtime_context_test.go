package goja

import (
	ct "context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSetContextCancelsGoCalls(t *testing.T) {
	r := New()
	ctx, cancel := ct.WithCancel(ct.Background())
	r.SetContext(ctx)
	r.SetContextValue("executionId", "test")

	var gotValue interface{}
	if err := r.Set("wait", func(ctx ct.Context) error {
		gotValue = ctx.Value("executionId")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			return errors.New("call outlived its runtime context")
		}
	}); err != nil {
		t.Fatal(err)
	}

	time.AfterFunc(50*time.Millisecond, cancel)
	_, err := r.RunString(`wait()`)
	if err == nil || !strings.Contains(err.Error(), ct.Canceled.Error()) {
		t.Fatalf("got %v, want the call cancelled with the runtime context", err)
	}
	if gotValue != "test" {
		t.Fatalf("context value = %v, want test", gotValue)
	}
}

func TestSetContextNilRestoresDefault(t *testing.T) {
	r := New()
	ctx, cancel := ct.WithCancel(ct.Background())
	cancel()
	r.SetContext(ctx)
	r.SetContext(nil)

	if err := r.Set("check", func(ctx ct.Context) error { return ctx.Err() }); err != nil {
		t.Fatal(err)
	}
	if _, err := r.RunString(`check()`); err != nil {
		t.Fatalf("call got a cancelled context after SetContext(nil): %v", err)
	}
}

func TestContextFuncMissingArgs(t *testing.T) {
	r := New()
	if err := r.Set("f", func(ctx ct.Context, a string, b int) string {
		return a + ":" + strconv.Itoa(b)
	}); err != nil {
		t.Fatal(err)
	}
	if err := r.Set("g", func(ctx ct.Context, a string, rest ...int) string {
		return a + ":" + strconv.Itoa(len(rest))
	}); err != nil {
		t.Fatal(err)
	}
	for src, want := range map[string]string{
		`f("x")`: "x:0",
		`f()`:    ":0",
		`g()`:    ":0",
		`g("y")`: "y:0",
	} {
		v, err := r.RunString(src)
		if err != nil {
			t.Fatalf("%s: %v", src, err)
		}
		if got := v.String(); got != want {
			t.Fatalf("%s = %q, want %q", src, got, want)
		}
	}
}
