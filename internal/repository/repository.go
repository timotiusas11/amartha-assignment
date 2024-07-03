package repository

import (
	"encoding/json"
	"fmt"

	"github.com/timotiusas11/amartha-assignment/common/driver/http"
	"github.com/timotiusas11/amartha-assignment/common/driver/inmemlib"
	"github.com/timotiusas11/amartha-assignment/common/driver/nsq"
	"github.com/timotiusas11/amartha-assignment/internal/model"
)

type RepositoryInterface interface {
	InsertLoan(loan model.Loan) error
	GetLoans() ([]model.Loan, error)
	GetLoan(loanID int64) (model.Loan, error)
	UpdateLoan(loan model.Loan) error
	Publish(loanID int64, invesment model.Investment) error
	GenerateAgreementLetter(loanID int64) error
}

type Repository struct {
	inmemlib   inmemlib.InMemLibInterface
	nsqClient  nsq.NSQInterface
	httpClient http.HTTPInterface
}

func NewRepository() Repository {
	return Repository{
		inmemlib:   inmemlib.New(),
		nsqClient:  nsq.New(),
		httpClient: http.New(),
	}
}

const (
	CacheKeyLoans = "loans"
)

func (r Repository) InsertLoan(loan model.Loan) error {
	var loanMap map[int64]model.Loan

	// Retrieve existing loan map from memcache
	exists, err := r.inmemlib.Get(CacheKeyLoans, func(val []byte) error {
		return json.Unmarshal(val, &loanMap)
	})

	if err != nil {
		return fmt.Errorf("failed to get loans from memcache: %w", err)
	}

	if !exists {
		// Cache miss, create a new map
		loanMap = make(map[int64]model.Loan)
	}

	// Add the new loan to the map
	loanMap[loan.LoanID] = loan

	// Set the updated map back into memcache
	err = r.inmemlib.Set(CacheKeyLoans, loanMap)
	if err != nil {
		return fmt.Errorf("failed to set updated loan data in memcache: %w", err)
	}

	return nil
}

func (r Repository) GetLoans() ([]model.Loan, error) {
	var loanMap map[int64]model.Loan

	// Retrieve existing loan map from memcache
	_, err := r.inmemlib.Get(CacheKeyLoans, func(val []byte) error {
		return json.Unmarshal(val, &loanMap)
	})

	if err != nil {
		return []model.Loan{}, fmt.Errorf("failed to get loans from memcache: %w", err)
	}

	var loans []model.Loan = make([]model.Loan, 0)

	for _, v := range loanMap {
		loans = append(loans, v)
	}

	return loans, nil
}

func (r Repository) GetLoan(loanID int64) (model.Loan, error) {
	var loanMap map[int64]model.Loan

	// Retrieve existing loan map from memcache
	_, err := r.inmemlib.Get(CacheKeyLoans, func(val []byte) error {
		return json.Unmarshal(val, &loanMap)
	})

	if err != nil {
		return model.Loan{}, fmt.Errorf("failed to get loans from memcache: %w", err)
	}

	return loanMap[loanID], nil
}

func (r Repository) UpdateLoan(loan model.Loan) error {
	loanMap := make(map[int64]model.Loan)

	// Retrieve existing loan map from memcache
	_, err := r.inmemlib.Get(CacheKeyLoans, func(val []byte) error {
		return json.Unmarshal(val, &loanMap)
	})
	if err != nil {
		return fmt.Errorf("failed to get loans from cache: %w", err)
	}

	// Update the loan in the map
	loanMap[loan.LoanID] = loan

	// Set the updated map back into memcache
	err = r.inmemlib.Set(CacheKeyLoans, loanMap)
	if err != nil {
		return fmt.Errorf("failed to update loans in cache: %w", err)
	}

	return nil
}

const (
	EmailAgreementLetterChannel    = "email_agreement_letter"
	GenerateAgreementLetterChannel = "generate_agreement_letter"
)

func (r Repository) Publish(loanID int64, invesment model.Investment) error {
	return r.nsqClient.Send(EmailAgreementLetterChannel, map[string]interface{}{
		"loan_id":         loanID,
		"investor_id":     invesment.InvestorID,
		"invested_amount": invesment.InvestedAmount,
	})
}

func (r Repository) GenerateAgreementLetter(loanID int64) error {
	return r.nsqClient.Send(GenerateAgreementLetterChannel, map[string]interface{}{
		"loan_id": loanID,
	})
}
