package commodityprovider

import (
	"atlas/internal/commoditydata"
	"context"
	"errors"
	"testing"
	"time"
)

type fake struct {
	name string
	o    commoditydata.Observation
	err  error
}

func (f fake) Name() string { return f.name }
func (f fake) Fetch(context.Context, QuoteRequest) (commoditydata.Observation, error) {
	return f.o, f.err
}
func TestMultiSourceIngestionAndResolution(t *testing.T) {
	fab := commoditydata.New()
	ing := New(fab)
	now := time.Now().UTC()
	mk := func(src string, v float64) commoditydata.Observation {
		return commoditydata.Observation{Commodity: "gold", Metric: "spot", Value: v, Unit: "oz", Currency: "USD", Source: src, ObservedAt: now, MaxAge: time.Hour, Confidence: .95, License: commoditydata.LicensePolicy{AllowDecisionUse: true}, Evidence: []string{"provider:" + src}}
	}
	ing.Register(fake{"a", mk("a", 2400), nil}, Config{Priority: 1})
	ing.Register(fake{"b", mk("b", 2401), nil}, Config{Priority: 2})
	r, e := ing.Fetch(context.Background(), QuoteRequest{Commodity: "gold", Metric: "spot"})
	if e != nil || len(r.Succeeded) != 2 {
		t.Fatalf("%v %#v", e, r)
	}
	if _, e = fab.Resolve(commoditydata.Request{Commodity: "gold", Metric: "spot", At: now.Add(time.Minute), RequireDecisionUse: true, MinConfidence: .9}); e != nil {
		t.Fatal(e)
	}
}
func TestProviderFailureDoesNotBlockHealthyFallback(t *testing.T) {
	fab := commoditydata.New()
	ing := New(fab)
	now := time.Now().UTC()
	ing.Register(fake{"bad", commoditydata.Observation{}, errors.New("down")}, Config{Priority: 1, MaxFailures: 1})
	ing.Register(fake{"good", commoditydata.Observation{Commodity: "wheat", Metric: "price", Value: 250, Source: "good", ObservedAt: now, MaxAge: time.Hour, Confidence: .9, License: commoditydata.LicensePolicy{AllowDecisionUse: true}, Evidence: []string{"signed-feed"}}, nil}, Config{Priority: 2})
	r, e := ing.Fetch(context.Background(), QuoteRequest{Commodity: "wheat", Metric: "price"})
	if e != nil || len(r.Succeeded) != 1 || r.Succeeded[0] != "good" {
		t.Fatalf("%v %#v", e, r)
	}
}
