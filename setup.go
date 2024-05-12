package stats

import (
	"strconv"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const PluginName = "stats"
const RequiredArgs = 1

func init() {
	plugin.Register(PluginName, setup)
}

func setup(c *caddy.Controller) error {
	var args []string
	c.NextArg() // Skip the name of the plugin, which is returned as an argument
	for c.NextArg() {
		args = append(args, c.Val())
	}

	if len(args) < RequiredArgs {
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error(PluginName, c.ArgErr())
	}

	backendURI := args[0]

	workers := int64(3)
	queryTimeout := 5 * time.Second
	statsPrefix := "coredns"
	maxEntryAge := 30 * 24 * time.Hour

	for c.NextBlock() {
		switch c.Val() {
		case "workers":
			t := c.RemainingArgs()
			if len(t) != 1 {
				return plugin.Error(PluginName, c.Errf("workers must be an integer"))
			}
			var err error
			workers, err = strconv.ParseInt(t[0], 10, 64)
			if err != nil {
				return plugin.Error(PluginName, c.Errf("workers must be an integer: %s", err.Error()))
			}
		case "queryTimeout":
			t := c.RemainingArgs()
			if len(t) != 1 {
				return plugin.Error(PluginName, c.Errf("queryTimeout must be an integer"))
			}
			var err error
			queryTimeout, err = time.ParseDuration(t[0])
			if err != nil {
				return plugin.Error(PluginName, c.Errf("queryTimeout must be a duration: %s", err.Error()))
			}
		case "statsPrefix":
			t := c.RemainingArgs()
			if len(t) != 1 {
				return plugin.Error(PluginName, c.Errf("statsPrefix must be a string"))
			}
			statsPrefix = t[0]
		case "maxEntryAge":
			t := c.RemainingArgs()
			if len(t) != 1 {
				return plugin.Error(PluginName, c.Errf("maxEntryAge must be a duration"))
			}
			var err error
			maxEntryAge, err = time.ParseDuration(t[0])
			if err != nil {
				return plugin.Error(PluginName, c.Errf("maxEntryAge must be a duration: %s", err.Error()))
			}

		default:
			return plugin.Error(PluginName, c.Errf("unknown property '%s'", c.Val()))
		}
	}

	logger := clog.NewWithPlugin(PluginName)

	backend, err := PrepareStatsBackend(backendURI, workers, queryTimeout, statsPrefix, maxEntryAge, logger)
	if err != nil {
		return plugin.Error(PluginName, err)
	}

	c.OnShutdown(backend.Stop)

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Stats{
			Next:    next,
			Backend: backend,

			Logger: logger,
		}
	})

	// All OK, return a nil error.
	return nil
}
