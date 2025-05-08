package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Get("/", app.Home)
	mux.Get("/receipt", app.Receipt)
	mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/venue/{id}", app.ChargeOnce)
	fileServer := http.FileServer(http.Dir("./static"))

	//auth routes
	mux.Get("/login", app.LoginPage)
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
