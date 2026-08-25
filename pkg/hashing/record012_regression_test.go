package hashing

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
)

func TestSHA256DoesNotCarryState(t *testing.T) {
	want := SHA256([]byte("component-a"))
	if got := SHA256([]byte("component-a")); got != want {
		t.Fatalf("repeated digest changed: got %x want %x", got, want)
	}
}

func TestHexConcurrentDeterministic(t *testing.T) {
	inputs := [][]byte{[]byte("alpha"), []byte("beta"), []byte("gamma"), []byte("delta")}
	want := make([]string, len(inputs))
	for i, input := range inputs {
		want[i] = Hex(input)
	}

	const rounds = 64
	var wg sync.WaitGroup
	start := make(chan struct{})
	for round := 0; round < rounds; round++ {
		for i, input := range inputs {
			wg.Add(1)
			go func(index int, data []byte) {
				defer wg.Done()
				<-start
				if got := Hex(data); got != want[index] {
					t.Errorf("digest mismatch: got %s want %s", got, want[index])
				}
			}(i, append([]byte(nil), input...))
		}
	}
	close(start)
	wg.Wait()
}

func TestEqualDoesNotRetainPriorMismatch(t *testing.T) {
	if Equal([]byte("first"), []byte("other")) {
		t.Fatal("different values compared equal")
	}
	if !Equal([]byte("stable"), []byte("stable")) {
		t.Fatal("an earlier mismatch contaminated a later comparison")
	}
}

func TestHexReaderDoesNotCarryState(t *testing.T) {
	want := Hex([]byte("document"))
	first, err := HexReader(strings.NewReader("document"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := HexReader(bytes.NewReader([]byte("document")))
	if err != nil {
		t.Fatal(err)
	}
	if first != want || second != want {
		t.Fatalf("reader digest leaked state: first=%s second=%s want=%s", first, second, want)
	}
}

type failingReader struct {
	delivered bool
	err       error
}

func (r *failingReader) Read(p []byte) (int, error) {
	if r.delivered {
		return 0, r.err
	}
	r.delivered = true
	return copy(p, "partial"), nil
}

func TestHexReaderPropagatesReadFailure(t *testing.T) {
	wantErr := errors.New("source interrupted")
	got, err := HexReader(&failingReader{err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected read error, got digest=%q err=%v", got, err)
	}
	if got != "" {
		t.Fatalf("partial digest was published: %q", got)
	}
}

var _ io.Reader = (*failingReader)(nil)
