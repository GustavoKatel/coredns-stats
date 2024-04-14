package stats

import (
	"fmt"
	"strings"
	"time"
)

type StatsBackend interface {
	Start() error
	Store(domain string, questionType string, metadata Metadata)
	Stop() error
}

func PrepareStatsBackend(uri string, workers int64, queryTimeout time.Duration, statsPrefix string, logger Logger) (StatsBackend, error) {
	if strings.HasPrefix(uri, "postgresql://") || strings.HasPrefix(uri, "postgres://") {
		backend := newBackendPostgres(uri, workers, queryTimeout, statsPrefix, logger)
		if err := backend.Start(); err != nil {
			return nil, err
		}

		return backend, nil
	}

	return nil, fmt.Errorf("unsupported backend: %s", uri)
}
