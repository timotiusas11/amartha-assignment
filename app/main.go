package main

import (
	"fmt"
	"net/http"

	"github.com/timotiusas11/amartha-assignment/internal/delivery"
	"github.com/timotiusas11/amartha-assignment/internal/repository"
	"github.com/timotiusas11/amartha-assignment/internal/usecase"
)

type application struct {
	router       *http.ServeMux
	deliveries   delivery.Delivery
	usecases     usecase.Usecase
	repositories repository.Repository
}

func newApplication() *application {
	return &application{
		router: http.NewServeMux(),
	}
}

func (a *application) repository() *application {
	a.repositories = repository.NewRepository()
	return a
}

func (a *application) usecase() *application {
	a.usecases = usecase.NewUsecase(a.repositories)
	return a
}

func (a *application) delivery() *application {
	a.deliveries = delivery.NewDelivery(a.usecases)

	// Health check
	a.router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Healthy"}`))
	})

	a.router.HandleFunc("/loans", a.deliveries.Loan)
	a.router.HandleFunc("/loans/{loan_id}", a.deliveries.GetLoan)
	a.router.HandleFunc("/loans/{loan_id}/approve", a.deliveries.Approve)
	a.router.HandleFunc("/loans/{loan_id}/invest", a.deliveries.Invest)
	a.router.HandleFunc("/loans/{loan_id}/disburse", a.deliveries.Disburse)

	// For admin only
	a.router.HandleFunc("/admin/view/loans", a.deliveries.AdminViewLoans)

	return a
}

func (a *application) serve() {
	fmt.Println("Server started at localhost:8080")
	http.ListenAndServe(":8080", a.router)
}

func main() {
	newApplication().repository().usecase().delivery().serve()
}
