package model

import "time"

type StateEnum int16

const (
	StateEnumProposed = iota
	StateEnumApproved
	StateEnumInvested
	StateEnumDisbursed
)

type Loan struct {
	LoanID             int64            `json:"loan_id"`              // Unique identifier
	BorrowerID         int64            `json:"borrower_id"`          // Identifier of the borrower
	PrincipalAmount    float64          `json:"principal_amount"`     // Amount of the loan requested
	Rate               float64          `json:"rate"`                 // Interest rate for the loan
	ROI                float64          `json:"roi"`                  // Return on investment for investors
	State              StateEnum        `json:"state"`                // Current state of the loan: proposed, approved, invested, disbursed
	ApprovalInfo       ApprovalInfo     `json:"approval_info"`        // Details when state is approved
	Investments        []Investment     `json:"investments"`          // List of investments and their invested amounts when state is invested
	DisbursementInfo   DisbursementInfo `json:"disbursement_info"`    // Details when state is disbursed
	AgreementLetterURL string           `json:"agreement_letter_url"` // Generated agreement letter
}

type ApprovalInfo struct {
	PictureProofURL  string    `json:"picture_proof_url"`
	FieldValidatorID int64     `json:"field_validator_id"`
	ApprovalDate     time.Time `json:"approval_date"`
}

type Investment struct {
	InvestorID     int64   `json:"investor_id"`
	InvestedAmount float64 `json:"invested_amount"`
}

type DisbursementInfo struct {
	SignedAgreementLetterURL string    `json:"signed_agreement_letter_url"`
	FieldOfficerID           int64     `json:"field_officer_id"`
	DisbursementDate         time.Time `json:"disbursement_date"`
}

type LoanInformation struct {
	LoanID             int64   `json:"loan_id"`
	BorrowerID         int64   `json:"borrower_id"`
	PrincipalAmount    float64 `json:"principal_amount"`
	Rate               float64 `json:"rate"`
	ROI                float64 `json:"roi"`
	AgreementLetterURL string  `json:"agreement_letter_url"`
}
