package trajectory

import (
	"errors"
	"time"
)

type Status string
const (
	Active Status = "ACTIVE"
	Held Status = "HELD"
	Terminated Status = "TERMINATED"
	Completed Status = "COMPLETED"
)

var (
	ErrNotActive = errors.New("trajectory not active")
	ErrBudgetExceeded = errors.New("trajectory budget exceeded")
	ErrRiskExceeded = errors.New("trajectory risk exceeded")
)

type State struct {
	TrajectoryID string
	Version uint64
	Status Status
	StartedAt time.Time
	ActionCount uint64
	RetryCount uint64
	SpendCents uint64
	TokenUsed uint64
	RiskScore uint32
	CompletedSteps map[string]bool
	Bindings map[string]string
}

type Action struct {
	ID string
	Name string
	Arguments map[string]string
	CostCents uint64
	Tokens uint64
	IsRetry bool
}

type Limits struct {
	MaxActions uint64
	MaxRetries uint64
	MaxSpendCents uint64
	MaxTokens uint64
	MaxRiskScore uint32
	MaxDuration time.Duration
}

func Evaluate(now time.Time,s State,a Action,l Limits) error {
	if s.Status != Active { return ErrNotActive }
	if now.Sub(s.StartedAt) > l.MaxDuration { return ErrBudgetExceeded }
	if s.ActionCount+1 > l.MaxActions { return ErrBudgetExceeded }
	if a.IsRetry && s.RetryCount+1 > l.MaxRetries { return ErrBudgetExceeded }
	if s.SpendCents+a.CostCents > l.MaxSpendCents { return ErrBudgetExceeded }
	if s.TokenUsed+a.Tokens > l.MaxTokens { return ErrBudgetExceeded }
	if s.RiskScore > l.MaxRiskScore { return ErrRiskExceeded }
	return nil
}
