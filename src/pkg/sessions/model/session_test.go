package model

import (
	"testing"
	"time"
)

func TestCloseSession(t *testing.T) {
	s := Session{}

	now := time.Now()
	s.Close()
	if !s.IsClosed() {
		t.Error("Session not closed")
	}
	if s.ClosedAt == nil {
		t.Error("ClosedAt not set")
	}
	if !s.ClosedAt.After(now) {
		t.Error("ClosedAt is less than time before close session")
	}
}
