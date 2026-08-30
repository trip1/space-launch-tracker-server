package notifications

import (
	"testing"

	"ds9labs.com/space-launch-server/internal/model"
)

func TestChangeDetectorSkipsFirstObservationAndEmitsNetChange(t *testing.T) {
	detector := NewChangeDetector()
	original := launch("l-1", "Falcon 9", "2026-09-01T12:00:00Z", 1)
	if got := detector.Changes([]model.Launch{original}); len(got) != 0 {
		t.Fatalf("first observation = %#v, want no notification", got)
	}
	if got := detector.Changes([]model.Launch{original}); len(got) != 0 {
		t.Fatalf("unchanged launch = %#v, want no notification", got)
	}

	moved := launch("l-1", "Falcon 9", "2026-09-01T14:00:00Z", 1)
	got := detector.Changes([]model.Launch{moved})
	if len(got) != 1 {
		t.Fatalf("changes = %#v, want one", got)
	}
	if got[0].Kind != LaunchChangeTime || got[0].LaunchID != "l-1" {
		t.Fatalf("change = %#v, want time change for l-1", got[0])
	}
}

func TestChangeDetectorEmitsStatusChange(t *testing.T) {
	detector := NewChangeDetector()
	detector.Changes([]model.Launch{launch("l-1", "Falcon 9", "2026-09-01T12:00:00Z", 1)})
	got := detector.Changes([]model.Launch{launch("l-1", "Falcon 9", "2026-09-01T12:00:00Z", 3)})
	if len(got) != 1 || got[0].Kind != LaunchChangeStatus {
		t.Fatalf("changes = %#v, want one status change", got)
	}
}

func launch(id, name, net string, statusID int) model.Launch {
	return model.Launch{ID: id, Name: name, Net: net, Status: model.LaunchStatus{ID: statusID}}
}
