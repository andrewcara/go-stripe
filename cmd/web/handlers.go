package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/andrewcara/go-stripe.git/internal/cards"
	"github.com/andrewcara/go-stripe.git/internal/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	eventID, _ := strconv.Atoi(r.Form.Get("product_id"))
	email := r.Form.Get("email")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	pm, err := card.GetPaymentMethod(paymentMethod)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	//create a new customer
	customerID, err := app.SaveCustomer(firstName, lastName, email)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println(customerID)
	//create new transaction
	amount, _ := strconv.Atoi(paymentMethod)

	txn := models.Transaction{
		Amount:              amount,
		Currency:            paymentCurrency,
		LastFour:            lastFour,
		ExpiryMonth:         int(expiryMonth),
		ExpiryYear:          int(expiryYear),
		BankReturnCode:      pi.Charges.Data[0].ID,
		TransactionStatusID: 2,
		PaymentIntent:       paymentIntent,
		PaymentMethod:       paymentMethod,
	}

	txnID, err := app.SaveTransaction(txn)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//create a new order

	order := models.Order{
		TicketID:      eventID,
		StatusID:      1,
		TransactionID: txnID,
		CustomerID:    customerID,
		Quantity:      1,
		Amount:        amount,
	}

	_, err = app.SaveOrder(order)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = pi.Charges.Data[0].ID
	data["first_name"] = firstName
	data["last_name"] = lastName

	//write this data to session then redirect so that we don't repost form data for payment
	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// saves a customer and returns an id
func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// save a transaction and return an ID
func (app *application) SaveTransaction(txn models.Transaction) (int, error) {

	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (app *application) SaveOrder(order models.Order) (int, error) {

	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	pk_key := os.Getenv("STRIPE_KEY")
	data := map[string]interface{}{
		"pk_key": pk_key,
	}

	id := chi.URLParam(r, "id")
	eventID, _ := strconv.Atoi(id)

	event, err := app.DB.GetTicketEvents(eventID)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data["event"] = event

	td := &templateData{Data: data}
	if err := app.renderTemplate(w, r, "buy-once", td, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
