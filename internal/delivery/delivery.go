package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/timotiusas11/amartha-assignment/internal/model"
	"github.com/timotiusas11/amartha-assignment/internal/usecase"
)

type Delivery struct {
	usecase.UsecaseInterface
}

func NewDelivery(usecase usecase.UsecaseInterface) Delivery {
	return Delivery{
		UsecaseInterface: usecase,
	}
}

func (d Delivery) Loan(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method == http.MethodPost {
		d.createLoan(w, r)
		return
	}

	// Check if the method is GET
	if r.Method == http.MethodGet {
		d.getLoans(w, r)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func (d Delivery) createLoan(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into LoanRequest struct
	var loan model.Loan
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Call the usecase's CreateLoan method
	err = d.UsecaseInterface.CreateLoan(loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.ROI)
	if err != nil {
		http.Error(w, "Failed to create loan", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Loan created successfully"))
}

func (d Delivery) getLoans(w http.ResponseWriter, _ *http.Request) {
	// Call the usecase's GetLoans method
	loans, err := d.UsecaseInterface.GetLoans()
	if err != nil {
		http.Error(w, "Failed to get loans", http.StatusInternalServerError)
		return
	}

	// Convert loans to JSON
	loansJSON, err := json.Marshal(loans)
	if err != nil {
		http.Error(w, "Failed to marshal loans data", http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(loansJSON)
}

func (d Delivery) GetLoan(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the loan ID from the URL path
	loanIDString := r.PathValue("loan_id")

	loanID, err := strconv.ParseInt(loanIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Call the usecase's GetLoan method
	loan, err := d.UsecaseInterface.GetLoan(loanID)
	if err != nil {
		http.Error(w, "Failed to get loan", http.StatusInternalServerError)
		return
	}

	// Send the loan details in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loan)
}

func (d Delivery) Approve(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the loan ID from the URL path
	loanIDString := r.PathValue("loan_id")
	loanID, err := strconv.ParseInt(loanIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Parse the incoming JSON request
	var approval model.ApprovalInfo
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&approval)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the usecase's Approve method
	err = d.UsecaseInterface.Approve(loanID, approval.PictureProofURL, approval.FieldValidatorID)
	if err != nil {
		http.Error(w, "Failed to approve loan", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan approved successfully"))
}

func (d Delivery) Invest(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the loan ID from the URL path
	loanIDString := r.PathValue("loan_id")
	loanID, err := strconv.ParseInt(loanIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Parse the incoming JSON request
	var invest model.Investment
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&invest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the usecase's Invest method
	err = d.UsecaseInterface.Invest(loanID, invest)
	if err != nil {
		http.Error(w, "Failed to invest", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Investment successful"))
}

func (d Delivery) Disburse(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the loan ID from the URL path
	loanIDString := r.PathValue("loan_id")
	loanID, err := strconv.ParseInt(loanIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var disbursementInfo model.DisbursementInfo
	err = json.NewDecoder(r.Body).Decode(&disbursementInfo)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Call usecase.Disburse
	err = d.UsecaseInterface.Disburse(loanID, disbursementInfo.SignedAgreementLetterURL, disbursementInfo.FieldOfficerID)
	if err != nil {
		http.Error(w, "Failed to disburse loan", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan disbursed successfully"))
}

func (d Delivery) AdminViewLoans(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Call the usecase's GetLoans method
	loans, err := d.UsecaseInterface.AdminViewLoans()
	if err != nil {
		http.Error(w, "Failed to get loans", http.StatusInternalServerError)
		return
	}

	// Convert loans to JSON
	loansJSON, err := json.Marshal(loans)
	if err != nil {
		http.Error(w, "Failed to marshal loans data", http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(loansJSON)
}
