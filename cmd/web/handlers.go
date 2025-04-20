package main

import (
	"net/http"
	"os"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	pk_key := os.Getenv("STRIPE_KEY")
	data := map[string]interface{}{
		"pk_key": pk_key,
	}
	td := &templateData{Data: data}

	if err := app.renderTemplate(w, r, "terminal", td, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	cardHolder := r.Form.Get("cardholder_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	data := make(map[string]interface{})
	data["cardholder_name"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	pk_key := os.Getenv("STRIPE_KEY")
	data := map[string]interface{}{
		"pk_key": pk_key,
	}
	td := &templateData{Data: data}
	if err := app.renderTemplate(w, r, "buy-once", td, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
