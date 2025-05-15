package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.Get("/", app.Home)
	mux.Get("/receipt", app.Receipt)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/virtual-terminal", app.VirtualTerminal)
	})
	// mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	// mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/venue/{id}", app.ChargeOnce)
	fileServer := http.FileServer(http.Dir("./static"))

	//auth routes
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
