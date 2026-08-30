package notifications

import "ds9labs.com/space-launch-server/internal/model"

type LaunchChangeKind string

const (
	LaunchChangeTime   LaunchChangeKind = "launch_time_changed"
	LaunchChangeStatus LaunchChangeKind = "launch_status_changed"
)

type LaunchChange struct {
	Kind     LaunchChangeKind
	LaunchID string
	Name     string
	Net      string
	StatusID int
}

type launchRevision struct {
	net      string
	statusID int
}

// ChangeDetector keeps only the fields that are safe and meaningful to use as
// notification triggers. The first observation establishes a baseline and
// deliberately emits nothing, preventing a restart from flooding devices.
type ChangeDetector struct {
	known map[string]launchRevision
}

func NewChangeDetector() *ChangeDetector {
	return &ChangeDetector{known: make(map[string]launchRevision)}
}

func (d *ChangeDetector) Changes(launches []model.Launch) []LaunchChange {
	changes := make([]LaunchChange, 0)
	for _, launch := range launches {
		if launch.ID == "" {
			continue
		}
		next := launchRevision{net: launch.Net, statusID: launch.Status.ID}
		previous, exists := d.known[launch.ID]
		d.known[launch.ID] = next
		if !exists {
			continue
		}
		if previous.statusID != next.statusID {
			changes = append(changes, LaunchChange{Kind: LaunchChangeStatus, LaunchID: launch.ID, Name: launch.Name, Net: launch.Net, StatusID: launch.Status.ID})
			continue
		}
		if previous.net != next.net {
			changes = append(changes, LaunchChange{Kind: LaunchChangeTime, LaunchID: launch.ID, Name: launch.Name, Net: launch.Net, StatusID: launch.Status.ID})
		}
	}
	return changes
}
