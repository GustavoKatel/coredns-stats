package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron"
)

type StatsBackend interface {
	Start() error
	Ready() bool
	Store(domain string, questionType string, metadata Metadata)
	Stop() error
}

func PrepareStatsBackend(uri string, workers int64, queryTimeout time.Duration, statsPrefix string, maxEntryAge time.Duration, maxEntryCleanCron string, logger Logger) (StatsBackend, error) {
	maxEntryCleanTicker := make(chan struct{}, 5)

	c := cron.New()
	if err := c.AddFunc(maxEntryCleanCron, func() {
		maxEntryCleanTicker <- struct{}{}
	}); err != nil {
		return nil, fmt.Errorf("error adding maxEntryCleanCron: %w", err)
	}

	if strings.HasPrefix(uri, "postgresql://") || strings.HasPrefix(uri, "postgres://") {
		backend := newBackendPostgres(uri, workers, queryTimeout, statsPrefix, maxEntryAge, maxEntryCleanTicker, logger)
		if err := backend.Start(); err != nil {
			return nil, err
		}

		c.Start()

		return backend, nil
	}

	return nil, fmt.Errorf("unsupported backend: %s", uri)
}
