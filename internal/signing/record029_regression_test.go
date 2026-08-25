package signing

import (
	"context"
	"errors"
	"testing"
)

func TestTypedNilPrimaryFallsBack(t *testing.T) {
	var primary *stubSigner
	service, err := NewService(primary, &stubSigner{value: "fallback-signature"})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Sign(context.Background(), "att-fallback", []byte("payload")); err != nil {
		t.Fatal(err)
	}
	if got := service.Report().Signatures["att-fallback"].Value; got != "fallback-signature" {
		t.Fatalf("fallback signature = %q", got)
	}
}

func TestSignerFailureRemainsClassifiable(t *testing.T) {
	service, err := NewService(&stubSigner{err: errKMSUnavailable}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = service.Sign(context.Background(), "att-failure", []byte("payload"))
	if !errors.Is(err, errKMSUnavailable) {
		t.Fatalf("signing error was not preserved: %v", err)
	}
}

func TestAbortDoesNotCommitPendingSignature(t *testing.T) {
	ledger := NewLedger()
	ledger.Stage("att-abort")
	ledger.Abort("att-abort")
	if signatures := ledger.snapshot(); len(signatures) != 0 {
		t.Fatalf("aborted signature was committed: %#v", signatures)
	}
}

func TestSuccessfulSigningReportIsDetached(t *testing.T) {
	service, err := NewService(&stubSigner{value: "signature"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Sign(context.Background(), "att-report", []byte("payload")); err != nil {
		t.Fatal(err)
	}
	report := service.Report()
	delete(report.Signatures, "att-report")
	if got := service.Report().Signatures["att-report"].Value; got != "signature" {
		t.Fatalf("report mutation changed signing history: %q", got)
	}
}
