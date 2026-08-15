package settlementgateway

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockProvider struct {
	mu            sync.Mutex
	sim           Simulation
	exec          Execution
	fin           Finality
	simErr        error
	execErr       error
	finErr        error
	executeCalls  int
	finalityCalls int
}

func (m *mockProvider) Simulate(context.Context, Intent) (Simulation, error) { return m.sim, m.simErr }
func (m *mockProvider) Execute(context.Context, Intent, Simulation) (Execution, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executeCalls++
	return m.exec, m.execErr
}
func (m *mockProvider) Finality(context.Context, Execution) (Finality, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.finalityCalls++
	return m.fin, m.finErr
}
func (m *mockProvider) calls() int { m.mu.Lock(); defer m.mu.Unlock(); return m.executeCalls }

func policy() Policy {
	return Policy{
		AllowedProviders:      map[string]bool{"chainlink": true},
		AllowedNetworks:       map[string]bool{"testnet": true},
		AllowedAssets:         map[string]bool{"USDC": true},
		AllowedCurrencies:     map[string]bool{"USD": true},
		AllowedCounterparties: map[string]bool{"supplier-1": true},
		MaxSingleMinor:        100_000,
		MaxAggregateMinor:     150_000,
		RequireApprovalAbove:  50_000,
		MinFinality:           2,
		RequireSimulation:     true,
		RequireVerifiedInputs: true,
	}
}

func intent() Intent {
	return Intent{
		ID:               "intent-1",
		PrincipalID:      "principal-1",
		AgentID:          "agent-1",
		TrajectoryID:     "trajectory-1",
		TransactionID:    "trade-1",
		Provider:         "chainlink",
		Network:          "testnet",
		Asset:            "USDC",
		Currency:         "USD",
		Counterparty:     "supplier-1",
		Purpose:          "inspection fee",
		AuthorityDigest:  "authority-1",
		ComplianceDigest: "compliance-1",
		PolicyDigest:     "policy-1",
		AmountMinor:      10_000,
		MaxSlippageBPS:   25,
		EvidenceDigests:  []string{"evidence-1"},
		ExpiresAt:        time.Now().UTC().Add(time.Hour),
	}
}

func provider() *mockProvider {
	return &mockProvider{
		sim:  Simulation{Approved: true, FeeMinor: 10, SlippageBPS: 5, EvidenceDigest: "sim-evidence"},
		exec: Execution{ProviderReference: "tx-1", ProviderReceipt: "provider-receipt", ExecutedAt: time.Unix(1_800_000_000, 0).UTC()},
		fin:  Finality{Confirmed: true, Confirmations: 3, Proof: "finality-proof"},
	}
}

func gateway(t *testing.T, p *mockProvider) *Gateway {
	t.Helper()
	g := New(policy())
	if err := g.Register("chainlink", p); err != nil {
		t.Fatal(err)
	}
	return g
}

func TestExecuteConfirmedAndReceiptVerifies(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	r, err := g.Execute(context.Background(), intent())
	if err != nil {
		t.Fatal(err)
	}
	if r.Receipt.Status != StatusConfirmed || !VerifyReceipt(r.Receipt) {
		t.Fatalf("bad receipt: %+v", r.Receipt)
	}
	if g.SpentMinor() != 10_000 || p.calls() != 1 {
		t.Fatal("ledger/provider mismatch")
	}
}

func TestIdempotentReplayDoesNotExecuteTwice(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	first, err := g.Execute(context.Background(), intent())
	if err != nil {
		t.Fatal(err)
	}
	second, err := g.Execute(context.Background(), intent())
	if err != nil {
		t.Fatal(err)
	}
	if !second.Idempotent || second.Receipt.ReceiptDigest != first.Receipt.ReceiptDigest || p.calls() != 1 || g.SpentMinor() != 10_000 {
		t.Fatal("idempotency failed")
	}
}

func TestPolicyAndEvidenceGates(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Intent)
		want   error
	}{
		{"evidence", func(i *Intent) { i.EvidenceDigests = nil }, ErrEvidence},
		{"provider", func(i *Intent) { i.Provider = "unknown" }, ErrProvider},
		{"network", func(i *Intent) { i.Network = "mainnet" }, ErrPolicy},
		{"asset", func(i *Intent) { i.Asset = "ETH" }, ErrPolicy},
		{"currency", func(i *Intent) { i.Currency = "EUR" }, ErrPolicy},
		{"counterparty", func(i *Intent) { i.Counterparty = "unknown" }, ErrPolicy},
		{"authority", func(i *Intent) { i.AuthorityDigest = "" }, ErrEvidence},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := provider()
			g := gateway(t, p)
			in := intent()
			tc.mutate(&in)
			_, err := g.Execute(context.Background(), in)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
			if p.calls() != 0 {
				t.Fatal("provider called on denied request")
			}
		})
	}
}

func TestApprovalAndAggregateLimits(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	in := intent()
	in.AmountMinor = 60_000
	if _, err := g.Execute(context.Background(), in); !errors.Is(err, ErrApproval) {
		t.Fatal(err)
	}
	in.ApprovalDigest = "approval-1"
	if _, err := g.Execute(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	in2 := intent()
	in2.ID = "intent-2"
	in2.AmountMinor = 100_000
	in2.ApprovalDigest = "approval-2"
	if _, err := g.Execute(context.Background(), in2); !errors.Is(err, ErrAggregate) {
		t.Fatalf("got %v", err)
	}
}

func TestSimulationFailureAndSlippageRollbackReservation(t *testing.T) {
	p := provider()
	p.sim.Approved = false
	g := gateway(t, p)
	if _, err := g.Execute(context.Background(), intent()); !errors.Is(err, ErrSimulation) {
		t.Fatal(err)
	}
	if g.SpentMinor() != 0 || p.calls() != 0 {
		t.Fatal("reservation not rolled back")
	}

	p = provider()
	p.sim.SlippageBPS = 30
	g = gateway(t, p)
	if _, err := g.Execute(context.Background(), intent()); !errors.Is(err, ErrSlippage) {
		t.Fatal(err)
	}
	if g.SpentMinor() != 0 || p.calls() != 0 {
		t.Fatal("slippage reservation not rolled back")
	}
}

func TestProviderExecutionFailureRollsBackReservation(t *testing.T) {
	p := provider()
	p.execErr = errors.New("down")
	g := gateway(t, p)
	if _, err := g.Execute(context.Background(), intent()); !errors.Is(err, ErrProvider) {
		t.Fatal(err)
	}
	if g.SpentMinor() != 0 {
		t.Fatal("execution failure consumed budget")
	}
}

func TestFinalityFailureRequiresReconciliationWithoutDoubleSpend(t *testing.T) {
	p := provider()
	p.fin = Finality{Confirmed: false, Confirmations: 0}
	g := gateway(t, p)
	r, err := g.Execute(context.Background(), intent())
	if !errors.Is(err, ErrReconcile) || r.Receipt.Status != StatusReconcile {
		t.Fatalf("%v %+v", err, r)
	}
	if g.SpentMinor() != 10_000 || p.calls() != 1 {
		t.Fatal("executed settlement must remain economically accounted")
	}

	p.mu.Lock()
	p.fin = Finality{Confirmed: true, Confirmations: 4, Proof: "later-proof"}
	p.mu.Unlock()
	reconciled, err := g.Reconcile(context.Background(), "intent-1")
	if err != nil || reconciled.Status != StatusConfirmed || !VerifyReceipt(reconciled) {
		t.Fatalf("%v %+v", err, reconciled)
	}
	if p.calls() != 1 || g.SpentMinor() != 10_000 {
		t.Fatal("reconcile re-executed settlement")
	}
}

func TestExpiredIntentRejected(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	in := intent()
	in.ExpiresAt = time.Now().UTC().Add(-time.Second)
	if _, err := g.Execute(context.Background(), in); !errors.Is(err, ErrExpired) {
		t.Fatal(err)
	}
	if p.calls() != 0 {
		t.Fatal("expired intent reached provider")
	}
}

func TestReceiptTamperDetected(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	r, err := g.Execute(context.Background(), intent())
	if err != nil {
		t.Fatal(err)
	}
	r.Receipt.AmountMinor++
	if VerifyReceipt(r.Receipt) {
		t.Fatal("tampered receipt verified")
	}
}

func TestConcurrentDuplicateExecutesAtMostOnce(t *testing.T) {
	p := provider()
	g := gateway(t, p)
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := g.Execute(context.Background(), intent()); results <- err }()
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil && !errors.Is(err, ErrInFlight) {
			t.Fatal(err)
		}
	}
	if p.calls() != 1 {
		t.Fatalf("execute calls=%d", p.calls())
	}
}
