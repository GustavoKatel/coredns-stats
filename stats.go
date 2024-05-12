package stats

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/plugin/ready"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var _ plugin.Handler = (*Stats)(nil)
var _ metadata.Provider = (*Stats)(nil)
var _ ready.Readiness = (*Stats)(nil)

type Stats struct {
	Next    plugin.Handler
	Backend StatsBackend

	Logger Logger
}

func (s Stats) Name() string { return PluginName }

func (s Stats) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		// If the query has no question, then we can't do anything with it here. So, I will just
		// pass it to the next plugin.
		return s.Next.ServeDNS(ctx, w, r)
	}

	question := r.Question[0]
	domain := question.Name
	questionType := dns.TypeToString[question.Qtype]

	n, err := s.Next.ServeDNS(ctx, w, r)

	valueFuncs := metadata.ValueFuncs(ctx)
	mtValues := Metadata{}

	for key, valueFunc := range valueFuncs {
		mtValues[key] = valueFunc()
	}

	s.Backend.Store(domain, questionType, mtValues)

	return n, err
}

// Metadata adds metadata to the context and returns a (potentially) new context.
// Note: this method should work quickly, because it is called for every request
// from the metadata plugin.
func (s Stats) Metadata(ctx context.Context, state request.Request) context.Context {

	return ctx
}

func (s Stats) Ready() bool {
	return s.Backend.Ready()
}
