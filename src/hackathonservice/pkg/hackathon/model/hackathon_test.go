package model

import (
	"testing"
	"time"
)

func TestCloseHackathon(t *testing.T) {
	s := Hackathon{}

	now := time.Now()
	s.Close()
	if !s.IsClosed() {
		t.Error("Hackathon not closed")
	}
	if s.ClosedAt == nil {
		t.Error("ClosedAt not set")
	}
	if !s.ClosedAt.After(now) {
		t.Error("ClosedAt is less than time before close session")
	}
}
