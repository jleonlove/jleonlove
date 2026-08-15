package transactionverification

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var ErrInvalidAssessment = errors.New("invalid verification assessment")

type Status string

const (
	StatusPass   Status = "PASS"
	StatusReview Status = "REVIEW"
	StatusBlock  Status = "BLOCK"
)

type Check struct {
	ID          string
	Passed      bool
	Required    bool
	EvidenceRef string
	Source      string
	VerifiedAt  time.Time
	ExpiresAt   time.Time
}

type Assessment struct {
	DealID                                                       string
	SubjectID                                                    string
	KYC, KYB, UBO, Mandate, Sanctions, PEP, AdverseMedia         Check
	ProofOfCommodity, Facility, DocumentIntegrity, BankOwnership Check
}

type Result struct {
	Status                         Status
	Score                          int
	Passed, Failed, Missing, Stale []string
}

func checks(a Assessment) []Check {
	return []Check{a.KYC, a.KYB, a.UBO, a.Mandate, a.Sanctions, a.PEP, a.AdverseMedia, a.ProofOfCommodity, a.Facility, a.DocumentIntegrity, a.BankOwnership}
}

func Evaluate(a Assessment, now time.Time) (Result, error) {
	if strings.TrimSpace(a.DealID) == "" || strings.TrimSpace(a.SubjectID) == "" {
		return Result{}, ErrInvalidAssessment
	}
	r := Result{Status: StatusPass}
	total := 0
	earned := 0
	for _, c := range checks(a) {
		if !c.Required {
			continue
		}
		total++
		id := strings.TrimSpace(c.ID)
		if id == "" {
			id = "unnamed_check"
		}
		if strings.TrimSpace(c.EvidenceRef) == "" || strings.TrimSpace(c.Source) == "" {
			r.Missing = append(r.Missing, id)
			continue
		}
		if !c.ExpiresAt.IsZero() && !now.Before(c.ExpiresAt) {
			r.Stale = append(r.Stale, id)
			continue
		}
		if !c.Passed {
			r.Failed = append(r.Failed, id)
			continue
		}
		earned++
		r.Passed = append(r.Passed, id)
	}
	if total > 0 {
		r.Score = earned * 100 / total
	}
	sort.Strings(r.Passed)
	sort.Strings(r.Failed)
	sort.Strings(r.Missing)
	sort.Strings(r.Stale)
	if len(r.Failed) > 0 {
		r.Status = StatusBlock
	} else if len(r.Missing) > 0 || len(r.Stale) > 0 {
		r.Status = StatusReview
	}
	return r, nil
}

func CanExecute(r Result) bool { return r.Status == StatusPass && r.Score == 100 }
