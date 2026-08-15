package settlementgateway

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrInput      = errors.New("invalid settlement intent")
	ErrExpired    = errors.New("settlement intent expired")
	ErrEvidence   = errors.New("verified input evidence required")
	ErrPolicy     = errors.New("settlement policy denied")
	ErrApproval   = errors.New("settlement approval required")
	ErrAggregate  = errors.New("aggregate settlement limit exceeded")
	ErrSimulation = errors.New("settlement simulation failed")
	ErrSlippage   = errors.New("settlement slippage exceeds limit")
	ErrProvider   = errors.New("settlement provider unavailable")
	ErrInFlight   = errors.New("settlement intent already in flight")
	ErrFinality   = errors.New("settlement finality not satisfied")
	ErrReconcile  = errors.New("settlement requires reconciliation")
	ErrReceipt    = errors.New("invalid settlement receipt")
)

type Status string

const (
	StatusConfirmed Status = "CONFIRMED"
	StatusReconcile Status = "RECONCILE_REQUIRED"
)

type Intent struct {
	ID               string
	PrincipalID      string
	AgentID          string
	TrajectoryID     string
	TransactionID    string
	Provider         string
	Network          string
	Asset            string
	Currency         string
	Counterparty     string
	Purpose          string
	AmountMinor      int64
	MaxSlippageBPS   int64
	ApprovalDigest   string
	AuthorityDigest  string
	ComplianceDigest string
	PolicyDigest     string
	EvidenceDigests  []string
	ExpiresAt        time.Time
}

type Policy struct {
	AllowedProviders      map[string]bool
	AllowedNetworks       map[string]bool
	AllowedAssets         map[string]bool
	AllowedCurrencies     map[string]bool
	AllowedCounterparties map[string]bool
	MaxSingleMinor        int64
	MaxAggregateMinor     int64
	RequireApprovalAbove  int64
	MinFinality           uint64
	RequireSimulation     bool
	RequireVerifiedInputs bool
}

type Simulation struct {
	Approved       bool
	FeeMinor       int64
	SlippageBPS    int64
	EvidenceDigest string
}

type Execution struct {
	ProviderReference string
	ProviderReceipt   string
	ExecutedAt        time.Time
}

type Finality struct {
	Confirmed     bool
	Confirmations uint64
	Proof         string
}

type Provider interface {
	Simulate(context.Context, Intent) (Simulation, error)
	Execute(context.Context, Intent, Simulation) (Execution, error)
	Finality(context.Context, Execution) (Finality, error)
}

type Receipt struct {
	IntentID           string    `json:"intent_id"`
	PrincipalID        string    `json:"principal_id"`
	AgentID            string    `json:"agent_id"`
	TrajectoryID       string    `json:"trajectory_id"`
	TransactionID      string    `json:"transaction_id"`
	Provider           string    `json:"provider"`
	Network            string    `json:"network"`
	Asset              string    `json:"asset"`
	Currency           string    `json:"currency"`
	Counterparty       string    `json:"counterparty"`
	Purpose            string    `json:"purpose"`
	AuthorityDigest    string    `json:"authority_digest"`
	ComplianceDigest   string    `json:"compliance_digest"`
	PolicyDigest       string    `json:"policy_digest"`
	AmountMinor        int64     `json:"amount_minor"`
	ProviderReference  string    `json:"provider_reference"`
	ProviderReceipt    string    `json:"provider_receipt"`
	SimulationEvidence string    `json:"simulation_evidence"`
	FinalityProof      string    `json:"finality_proof"`
	Confirmations      uint64    `json:"confirmations"`
	Status             Status    `json:"status"`
	ExecutedAt         time.Time `json:"executed_at"`
	ReceiptDigest      string    `json:"receipt_digest"`
}

type Result struct {
	Receipt    Receipt
	Idempotent bool
}

type Gateway struct {
	mu        sync.Mutex
	policy    Policy
	providers map[string]Provider
	spent     int64
	inflight  map[string]struct{}
	receipts  map[string]Receipt
}

func New(policy Policy) *Gateway {
	return &Gateway{
		policy:    clonePolicy(policy),
		providers: map[string]Provider{},
		inflight:  map[string]struct{}{},
		receipts:  map[string]Receipt{},
	}
}

func (g *Gateway) Register(name string, p Provider) error {
	if strings.TrimSpace(name) == "" || p == nil {
		return ErrProvider
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.providers[name] = p
	return nil
}

func (g *Gateway) Execute(ctx context.Context, in Intent) (Result, error) {
	if err := validateIntent(in); err != nil {
		return Result{}, err
	}
	if !in.ExpiresAt.IsZero() && !time.Now().UTC().Before(in.ExpiresAt.UTC()) {
		return Result{}, ErrExpired
	}

	g.mu.Lock()
	if r, ok := g.receipts[in.ID]; ok {
		g.mu.Unlock()
		return Result{Receipt: r, Idempotent: true}, nil
	}
	if _, ok := g.inflight[in.ID]; ok {
		g.mu.Unlock()
		return Result{}, ErrInFlight
	}
	p, ok := g.providers[in.Provider]
	if !ok {
		g.mu.Unlock()
		return Result{}, ErrProvider
	}
	if err := authorize(g.policy, g.spent, in); err != nil {
		g.mu.Unlock()
		return Result{}, err
	}
	g.inflight[in.ID] = struct{}{}
	g.spent += in.AmountMinor // reserve before external execution; never oversubscribe concurrently.
	g.mu.Unlock()

	sim := Simulation{Approved: true}
	if g.policy.RequireSimulation {
		var err error
		sim, err = p.Simulate(ctx, in)
		if err != nil || !sim.Approved || strings.TrimSpace(sim.EvidenceDigest) == "" {
			g.rollbackReservation(in.ID, in.AmountMinor)
			return Result{}, ErrSimulation
		}
		if sim.SlippageBPS < 0 || (in.MaxSlippageBPS >= 0 && sim.SlippageBPS > in.MaxSlippageBPS) {
			g.rollbackReservation(in.ID, in.AmountMinor)
			return Result{}, ErrSlippage
		}
	}

	exec, err := p.Execute(ctx, in, sim)
	if err != nil || strings.TrimSpace(exec.ProviderReference) == "" || strings.TrimSpace(exec.ProviderReceipt) == "" {
		g.rollbackReservation(in.ID, in.AmountMinor)
		return Result{}, ErrProvider
	}
	if exec.ExecutedAt.IsZero() {
		exec.ExecutedAt = time.Now().UTC()
	}

	finality, ferr := p.Finality(ctx, exec)
	r := buildReceipt(in, sim, exec, finality)
	if ferr != nil || !finality.Confirmed || finality.Confirmations < g.policy.MinFinality || strings.TrimSpace(finality.Proof) == "" {
		r.Status = StatusReconcile
		r.ReceiptDigest = digestReceipt(r)
		g.finish(in.ID, r)
		return Result{Receipt: r}, ErrReconcile
	}
	r.Status = StatusConfirmed
	r.ReceiptDigest = digestReceipt(r)
	g.finish(in.ID, r)
	return Result{Receipt: r}, nil
}

// Reconcile re-checks finality without re-executing the economic action.
func (g *Gateway) Reconcile(ctx context.Context, intentID string) (Receipt, error) {
	g.mu.Lock()
	r, ok := g.receipts[intentID]
	if !ok || r.Status != StatusReconcile {
		g.mu.Unlock()
		return Receipt{}, ErrFinality
	}
	p, ok := g.providers[r.Provider]
	g.mu.Unlock()
	if !ok {
		return Receipt{}, ErrProvider
	}
	fin, err := p.Finality(ctx, Execution{ProviderReference: r.ProviderReference, ProviderReceipt: r.ProviderReceipt, ExecutedAt: r.ExecutedAt})
	if err != nil || !fin.Confirmed || fin.Confirmations < g.policy.MinFinality || strings.TrimSpace(fin.Proof) == "" {
		return r, ErrReconcile
	}
	r.Status = StatusConfirmed
	r.Confirmations = fin.Confirmations
	r.FinalityProof = fin.Proof
	r.ReceiptDigest = digestReceipt(r)
	g.mu.Lock()
	g.receipts[intentID] = r
	g.mu.Unlock()
	return r, nil
}

func (g *Gateway) Receipt(intentID string) (Receipt, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.receipts[intentID]
	return r, ok
}

func (g *Gateway) SpentMinor() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.spent
}

func VerifyReceipt(r Receipt) bool {
	if r.IntentID == "" || r.TransactionID == "" || r.ProviderReference == "" || r.ProviderReceipt == "" || r.AmountMinor <= 0 || r.ReceiptDigest == "" {
		return false
	}
	return digestReceipt(r) == r.ReceiptDigest
}

func authorize(p Policy, spent int64, in Intent) error {
	if p.RequireVerifiedInputs && len(in.EvidenceDigests) == 0 {
		return ErrEvidence
	}
	for _, d := range in.EvidenceDigests {
		if strings.TrimSpace(d) == "" {
			return ErrEvidence
		}
	}
	if strings.TrimSpace(in.AuthorityDigest) == "" || strings.TrimSpace(in.ComplianceDigest) == "" || strings.TrimSpace(in.PolicyDigest) == "" {
		return ErrEvidence
	}
	if !p.AllowedProviders[in.Provider] || !p.AllowedNetworks[in.Network] || !p.AllowedAssets[in.Asset] || !p.AllowedCurrencies[in.Currency] || !p.AllowedCounterparties[in.Counterparty] {
		return ErrPolicy
	}
	if p.MaxSingleMinor <= 0 || p.MaxAggregateMinor <= 0 || in.AmountMinor > p.MaxSingleMinor {
		return ErrPolicy
	}
	if spent > p.MaxAggregateMinor-in.AmountMinor {
		return ErrAggregate
	}
	if p.RequireApprovalAbove >= 0 && in.AmountMinor > p.RequireApprovalAbove && strings.TrimSpace(in.ApprovalDigest) == "" {
		return ErrApproval
	}
	return nil
}

func validateIntent(in Intent) error {
	if strings.TrimSpace(in.ID) == "" || strings.TrimSpace(in.PrincipalID) == "" || strings.TrimSpace(in.AgentID) == "" || strings.TrimSpace(in.TrajectoryID) == "" || strings.TrimSpace(in.TransactionID) == "" || strings.TrimSpace(in.Provider) == "" || strings.TrimSpace(in.Network) == "" || strings.TrimSpace(in.Asset) == "" || strings.TrimSpace(in.Currency) == "" || strings.TrimSpace(in.Counterparty) == "" || strings.TrimSpace(in.Purpose) == "" || in.AmountMinor <= 0 || in.MaxSlippageBPS < 0 {
		return ErrInput
	}
	return nil
}

func (g *Gateway) rollbackReservation(id string, amount int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.inflight, id)
	g.spent -= amount
	if g.spent < 0 {
		g.spent = 0
	}
}

func (g *Gateway) finish(id string, r Receipt) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.inflight, id)
	g.receipts[id] = r
}

func buildReceipt(in Intent, sim Simulation, exec Execution, fin Finality) Receipt {
	return Receipt{
		IntentID:           in.ID,
		PrincipalID:        in.PrincipalID,
		AgentID:            in.AgentID,
		TrajectoryID:       in.TrajectoryID,
		TransactionID:      in.TransactionID,
		Provider:           in.Provider,
		Network:            in.Network,
		Asset:              in.Asset,
		Currency:           in.Currency,
		Counterparty:       in.Counterparty,
		Purpose:            in.Purpose,
		AuthorityDigest:    in.AuthorityDigest,
		ComplianceDigest:   in.ComplianceDigest,
		PolicyDigest:       in.PolicyDigest,
		AmountMinor:        in.AmountMinor,
		ProviderReference:  exec.ProviderReference,
		ProviderReceipt:    exec.ProviderReceipt,
		SimulationEvidence: sim.EvidenceDigest,
		FinalityProof:      fin.Proof,
		Confirmations:      fin.Confirmations,
		ExecutedAt:         exec.ExecutedAt.UTC(),
	}
}

func digestReceipt(r Receipt) string {
	x := append([]string(nil), rFields(r)...)
	sort.Strings(x)
	h := sha256.Sum256([]byte(strings.Join(x, "\x00")))
	return hex.EncodeToString(h[:])
}

func rFields(r Receipt) []string {
	return []string{
		"intent=" + r.IntentID,
		"principal=" + r.PrincipalID,
		"agent=" + r.AgentID,
		"trajectory=" + r.TrajectoryID,
		"transaction=" + r.TransactionID,
		"provider=" + r.Provider,
		"network=" + r.Network,
		"asset=" + r.Asset,
		"currency=" + r.Currency,
		"counterparty=" + r.Counterparty,
		"purpose=" + r.Purpose,
		"authority=" + r.AuthorityDigest,
		"compliance=" + r.ComplianceDigest,
		"policy=" + r.PolicyDigest,
		"amount=" + itoa(r.AmountMinor),
		"provider_ref=" + r.ProviderReference,
		"provider_receipt=" + r.ProviderReceipt,
		"simulation=" + r.SimulationEvidence,
		"finality=" + r.FinalityProof,
		"confirmations=" + utoa(r.Confirmations),
		"status=" + string(r.Status),
		"executed_at=" + r.ExecutedAt.UTC().Format(time.RFC3339Nano),
	}
}

func clonePolicy(p Policy) Policy {
	p.AllowedProviders = cloneSet(p.AllowedProviders)
	p.AllowedNetworks = cloneSet(p.AllowedNetworks)
	p.AllowedAssets = cloneSet(p.AllowedAssets)
	p.AllowedCurrencies = cloneSet(p.AllowedCurrencies)
	p.AllowedCounterparties = cloneSet(p.AllowedCounterparties)
	return p
}

func cloneSet(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func utoa(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	return string(b[i:])
}
