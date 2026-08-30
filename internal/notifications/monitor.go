package notifications

import (
	"context"
	"log/slog"
	"time"

	"ds9labs.com/space-launch-server/internal/model"
)

type LaunchFeed interface {
	UpcomingLaunches(context.Context, int) (model.UpcomingLaunches, error)
}

type Monitor struct {
	feed     LaunchFeed
	registry *Registry
	sender   Sender
	detector *ChangeDetector
	logger   *slog.Logger
	limit    int
}

func NewMonitor(feed LaunchFeed, registry *Registry, sender Sender, logger *slog.Logger, limit int) *Monitor {
	return &Monitor{feed: feed, registry: registry, sender: sender, detector: NewChangeDetector(), logger: logger, limit: limit}
}
func (m *Monitor) Poll(ctx context.Context) {
	launches, err := m.feed.UpcomingLaunches(ctx, m.limit)
	if err != nil {
		m.logger.Warn("notification feed poll failed", "error", err)
		return
	}
	for _, change := range m.detector.Changes(launches.Items) {
		devices, err := m.registry.Devices(ctx)
		if err != nil {
			m.logger.Warn("load notification devices failed", "error", err)
			return
		}
		for _, device := range devices {
			err := m.sender.Send(ctx, device.Token, map[string]string{"event": "launch_changed", "launch_id": change.LaunchID, "kind": string(change.Kind)})
			if err != nil {
				m.logger.Warn("FCM delivery failed", "installation_id", device.InstallationID, "error", err)
			}
		}
	}
}
func (m *Monitor) Run(ctx context.Context, interval time.Duration) {
	m.Poll(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.Poll(ctx)
		}
	}
}
