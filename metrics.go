package stats

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// oldEntriesDeletedCount is number of old entries deleted in each maxEntryCleanCron tick
	oldEntriesDeletedCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: PluginName,
		Name:      "old_entries_deleted_count",
		Help:      "Number of old entries deleted in each maxEntryCleanCron tick",
	})
)
