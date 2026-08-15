package inspectionintel

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var ErrInvalidRequest = errors.New("invalid inspection request")
var ErrNoEligibleEvidence = errors.New("no eligible inspection evidence")
var ErrQualityFailed = errors.New("commodity quality failed")

type Evidence struct {
	ID, Provider, Commodity, LotID, FacilityID, Method, CertificateHash string
	ObservedAt, ValidUntil                                              time.Time
	Purity, Quantity, Confidence                                        float64
	DecisionUseAllowed, Independent, SignatureVerified                  bool
	MaxAge                                                              time.Duration
}
type Request struct {
	Commodity, LotID, FacilityID          string
	MinPurity, MinQuantity, MinConfidence float64
	At                                    time.Time
	RequireIndependent, RequireSignature  bool
}
type Decision struct {
	SelectedID, Provider, CertificateHash string
	Purity, Quantity, Confidence          float64
	Rejected                              map[string]string
}

func Evaluate(r Request, evidence []Evidence) (Decision, error) {
	if strings.TrimSpace(r.Commodity) == "" || strings.TrimSpace(r.LotID) == "" || r.MinPurity < 0 || r.MinQuantity < 0 {
		return Decision{}, ErrInvalidRequest
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	if r.MinConfidence <= 0 {
		r.MinConfidence = .8
	}
	d := Decision{Rejected: map[string]string{}}
	eligible := []Evidence{}
	for _, e := range evidence {
		key := e.ID
		if key == "" {
			key = e.Provider
		}
		why := ""
		switch {
		case e.ID == "" || e.Provider == "" || e.ObservedAt.IsZero() || e.Confidence < 0 || e.Confidence > 1 || e.Purity < 0 || e.Quantity < 0:
			why = "invalid"
		case !strings.EqualFold(e.Commodity, r.Commodity) || !strings.EqualFold(e.LotID, r.LotID):
			why = "scope_mismatch"
		case r.FacilityID != "" && !strings.EqualFold(e.FacilityID, r.FacilityID):
			why = "facility_mismatch"
		case !e.DecisionUseAllowed:
			why = "license_blocked"
		case r.RequireIndependent && !e.Independent:
			why = "not_independent"
		case r.RequireSignature && !e.SignatureVerified:
			why = "signature_unverified"
		case strings.TrimSpace(e.CertificateHash) == "":
			why = "missing_certificate_hash"
		case e.Confidence < r.MinConfidence:
			why = "low_confidence"
		case e.MaxAge <= 0 || r.At.Sub(e.ObservedAt) > e.MaxAge || e.ObservedAt.After(r.At.Add(time.Minute)):
			why = "stale_or_future"
		case !e.ValidUntil.IsZero() && r.At.After(e.ValidUntil):
			why = "expired"
		}
		if why != "" {
			d.Rejected[key] = why
			continue
		}
		eligible = append(eligible, e)
	}
	if len(eligible) == 0 {
		return d, ErrNoEligibleEvidence
	}
	sort.SliceStable(eligible, func(i, j int) bool {
		if eligible[i].Confidence == eligible[j].Confidence {
			return eligible[i].ID < eligible[j].ID
		}
		return eligible[i].Confidence > eligible[j].Confidence
	})
	e := eligible[0]
	d.SelectedID, d.Provider, d.CertificateHash = e.ID, e.Provider, e.CertificateHash
	d.Purity, d.Quantity, d.Confidence = e.Purity, e.Quantity, e.Confidence
	if e.Purity < r.MinPurity || e.Quantity < r.MinQuantity {
		return d, ErrQualityFailed
	}
	return d, nil
}
