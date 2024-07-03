package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/timotiusas11/amartha-assignment/internal/model"
	"github.com/timotiusas11/amartha-assignment/internal/repository"
)

type UsecaseInterface interface {
	CreateLoan(borrowerID int64, principalAmount float64, rate float64, roi float64) error
	GetLoans() ([]model.LoanInformation, error)
	GetLoan(loanID int64) (model.LoanInformation, error)
	Approve(loanID int64, pictureProofURL string, fieldValidatorID int64) error
	Invest(loanID int64, investment model.Investment) error
	Disburse(loanID int64, agreementLetterURL string, fieldOfficerID int64) error
	AdminViewLoans() ([]model.Loan, error)
}

type Usecase struct {
	repository.RepositoryInterface
}

func NewUsecase(repository repository.RepositoryInterface) Usecase {
	return Usecase{
		RepositoryInterface: repository,
	}
}

func (u Usecase) CreateLoan(borrowerID int64, principalAmount float64, rate float64, roi float64) error {
	// Create a new loan object
	loan := model.Loan{
		LoanID:          time.Now().UnixMilli(),
		BorrowerID:      borrowerID,
		PrincipalAmount: principalAmount,
		Rate:            rate,
		ROI:             roi,
		State:           model.StateEnumProposed,
	}

	// Call the dependency's InsertLoan method
	err := u.RepositoryInterface.InsertLoan(loan)
	if err != nil {
		return errors.New("failed to insert loan: " + err.Error())
	}

	// Generate agreement letter
	err = u.RepositoryInterface.GenerateAgreementLetter(loan.LoanID)
	if err != nil {
		return fmt.Errorf("failed to generate agreement letter: %w", err)
	}

	return nil
}

func (u Usecase) GetLoans() ([]model.LoanInformation, error) {
	// Call the repository's GetLoans method
	loans, err := u.RepositoryInterface.GetLoans()
	if err != nil {
		return nil, fmt.Errorf("failed to get loans from repository: %w", err)
	}

	var loanInformations []model.LoanInformation = make([]model.LoanInformation, 0)

	for _, loan := range loans {
		loanInformations = append(loanInformations, model.LoanInformation{
			LoanID:             loan.LoanID,
			BorrowerID:         loan.BorrowerID,
			PrincipalAmount:    loan.PrincipalAmount,
			Rate:               loan.Rate,
			ROI:                loan.ROI,
			AgreementLetterURL: loan.AgreementLetterURL,
		})
	}

	return loanInformations, nil
}

func (u Usecase) GetLoan(loanID int64) (model.LoanInformation, error) {
	// Call the repository's GetLoan method
	loan, err := u.RepositoryInterface.GetLoan(loanID)
	if err != nil {
		return model.LoanInformation{}, fmt.Errorf("failed to get loan from repository: %w", err)
	}

	// Return if the loan is not found
	if loan.LoanID == 0 {
		return model.LoanInformation{}, errors.New("loan not found")
	}

	return model.LoanInformation{
		LoanID:             loan.LoanID,
		BorrowerID:         loan.BorrowerID,
		PrincipalAmount:    loan.PrincipalAmount,
		Rate:               loan.Rate,
		ROI:                loan.ROI,
		AgreementLetterURL: loan.AgreementLetterURL,
	}, nil
}

func (u Usecase) Approve(loanID int64, pictureProofURL string, fieldValidatorID int64) error {
	// Check if any of the approval info fields are empty
	if pictureProofURL == "" || fieldValidatorID == 0 {
		return errors.New("approval info is incomplete")
	}

	// Retrieve the loan from the repository
	loan, err := u.RepositoryInterface.GetLoan(loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Return if the loan is not found
	if loan.LoanID == 0 {
		return errors.New("loan not found")
	}

	// Reject if the state of the loan is beyond approved
	if loan.State != model.StateEnumProposed {
		return errors.New("loan is already approved, invested, or disbursed")
	}

	// Update the loan's approval info and state
	loan.ApprovalInfo = model.ApprovalInfo{
		PictureProofURL:  pictureProofURL,
		FieldValidatorID: fieldValidatorID,
		ApprovalDate:     time.Now(),
	}
	loan.State = model.StateEnumApproved

	// Update the loan in the repository
	err = u.RepositoryInterface.UpdateLoan(loan)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	// Return success
	return nil
}

func (u Usecase) Invest(loanID int64, investment model.Investment) error {
	// Validate that the investment details are complete
	if investment.InvestorID == 0 || investment.InvestedAmount <= 0 {
		return errors.New("invalid investment details")
	}

	// Retrieve the loan from the repository
	loan, err := u.RepositoryInterface.GetLoan(loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Return if the loan is not found
	if loan.LoanID == 0 {
		return errors.New("loan not found")
	}

	// Ensure the loan can only be invested in if its state is approved
	if loan.State != model.StateEnumApproved {
		return errors.New("loan is not in approved state")
	}

	// Calculate the total invested amount
	totalInvestedAmount := investment.InvestedAmount
	for _, inv := range loan.Investments {
		totalInvestedAmount += inv.InvestedAmount
	}

	// Ensure the total invested amount does not exceed the loan principal amount
	if totalInvestedAmount > loan.PrincipalAmount {
		return errors.New("total invested amount exceeds principal amount")
	}

	// Update the investments of the loan
	loan.Investments = append(loan.Investments, investment)

	// If the total invested amount matches the principal amount, update the loan's status
	if totalInvestedAmount == loan.PrincipalAmount {
		loan.State = model.StateEnumInvested

		// Assume that the agreement letter has been generated in another service (background service)
		loan.AgreementLetterURL = "https://example.com/agreement_letter.pdf"

		// Send agreement letters to investors using a message queue service (NSQ)
		for _, inv := range loan.Investments {
			err = u.RepositoryInterface.Publish(loan.LoanID, inv)
			if err != nil {
				return fmt.Errorf("failed to publish agreement letter: %w", err)
			}
		}
	}

	// Update the loan in the cache
	err = u.RepositoryInterface.UpdateLoan(loan)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	return nil
}

func (u Usecase) Disburse(loanID int64, signedAgreementLetterURL string, fieldOfficerID int64) error {
	// Check if agreement letter URL or field officer ID is empty
	if signedAgreementLetterURL == "" || fieldOfficerID == 0 {
		return errors.New("agreement letter URL or field officer ID is empty")
	}

	// Retrieve the loan from the repository
	loan, err := u.RepositoryInterface.GetLoan(loanID)
	if err != nil {
		return fmt.Errorf("failed to get loan: %w", err)
	}

	// Return error if loan not found
	if loan.LoanID == 0 {
		return errors.New("loan not found")
	}

	// Loan can only be disbursed if its status is StateEnumInvested
	if loan.State != model.StateEnumInvested {
		return errors.New("loan can only be disbursed if status is invested")
	}

	// Update status of loan to StateEnumDisbursed
	loan.State = model.StateEnumDisbursed

	// Update disbursement info of the loan
	loan.DisbursementInfo = model.DisbursementInfo{
		SignedAgreementLetterURL: signedAgreementLetterURL,
		FieldOfficerID:           fieldOfficerID,
		DisbursementDate:         time.Now(),
	}

	// Update loan in the cache
	err = u.RepositoryInterface.UpdateLoan(loan)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	return nil
}

func (u Usecase) AdminViewLoans() ([]model.Loan, error) {
	// Call the repository's GetLoans method
	loans, err := u.RepositoryInterface.GetLoans()
	if err != nil {
		return nil, fmt.Errorf("failed to get loans from repository: %w", err)
	}

	return loans, nil
}
