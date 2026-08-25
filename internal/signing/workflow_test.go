package signing

import (
	"context"
	"errors"
	"testing"
)

var errKMSUnavailable = errors.New("kms unavailable")

type stubSigner struct {
	value string
	err   error
}

func (s *stubSigner) Sign(context.Context, []byte) (string, error) {
	if s == nil {
		panic("nil signer invoked")
	}
	return s.value, s.err
}

func TestSigningWorkflow(t *testing.T) {
	t.Run("typed nil uses fallback", func(t *testing.T) {
		var primary *stubSigner
		service, err := NewService(primary, &stubSigner{value: "fallback-signature"})
		if err != nil {
			t.Fatal(err)
		}
		if err := service.Sign(context.Background(), "att-1", []byte("payload")); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("failure is returned and rolled back", func(t *testing.T) {
		service, err := NewService(&stubSigner{err: errKMSUnavailable}, nil)
		if err != nil {
			t.Fatal(err)
		}
		err = service.Sign(context.Background(), "att-2", []byte("payload"))
		if !errors.Is(err, errKMSUnavailable) {
			t.Fatalf("signing error was not preserved: %v", err)
		}
		if signatures := service.Report().Signatures; len(signatures) != 0 {
			t.Fatalf("failed signature was committed: %#v", signatures)
		}
	})

	t.Run("successful report is detached", func(t *testing.T) {
		service, err := NewService(&stubSigner{value: "signature"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := service.Sign(context.Background(), "att-3", []byte("payload")); err != nil {
			t.Fatal(err)
		}
		report := service.Report()
		delete(report.Signatures, "att-3")
		if _, ok := service.Report().Signatures["att-3"]; !ok {
			t.Fatal("report mutation changed signing history")
		}
	})
}
