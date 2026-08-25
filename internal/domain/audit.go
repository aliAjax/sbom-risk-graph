package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

var lastAuditHash string

type AuditEvent struct {
	ID           string    `json:"id"`
	Action       string    `json:"action"`
	Subject      string    `json:"subject"`
	At           time.Time `json:"at"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
}

func (e AuditEvent) Seal() AuditEvent {
	if e.PreviousHash == "" && lastAuditHash != "" {
		e.PreviousHash = lastAuditHash
	}
	payload, _ := json.Marshal(struct {
		ID, Action, Subject, PreviousHash string
		At                                time.Time
	}{e.ID, e.Action, e.Subject, e.PreviousHash, e.At})
	d := sha256.Sum256(payload)
	e.Hash = hex.EncodeToString(d[:])
	lastAuditHash = e.Hash
	return e
}
