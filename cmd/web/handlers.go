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

	if err := app.renderTemplate(w, r, "terminal", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

//Get transaction data from post

func (app *application) GetTransactionData(r *http.Request) (TransactionData, error) {
	var txnData TransactionData
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	email := r.Form.Get("email")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	amount, _ := strconv.Atoi(paymentAmount)

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)

	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)

	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth),
		ExpiryYear:      int(expiryYear),
		BankReturnCode:  pi.Charges.Data[0].ID,
	}
	return txnData, nil
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	eventID, _ := strconv.Atoi(r.Form.Get("product_id"))

	txnData, err := app.GetTransactionData(r)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//create a new customer
	customerID, err := app.SaveCustomer(txnData.FirstName, txnData.LastName, txnData.Email)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
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
		Amount:        txnData.PaymentAmount,
	}

	_, err = app.SaveOrder(order)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//write this data to session then redirect so that we don't repost form data for payment

	app.Session.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/receipt", http.StatusSeeOther)

}
func (app *application) Receipt(w http.ResponseWriter, r *http.Request) {

	txn := app.Session.Get(r.Context(), "receipt").(TransactionData)
	data := make(map[string]interface{})
	data["txn"] = txn
	app.Session.Remove(r.Context(), "receipt")
	if err := app.renderTemplate(w, r, "receipt", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}
func (app *application) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {

	txnData, err := app.GetTransactionData(r)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
	}

	_, err = app.SaveTransaction(txn)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.Session.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/virtual-terminal-receipt", http.StatusSeeOther)

}
func (app *application) VirtualTerminalReceipt(w http.ResponseWriter, r *http.Request) {

	txn := app.Session.Get(r.Context(), "receipt").(TransactionData)
	data := make(map[string]interface{})
	data["txn"] = txn
	app.Session.Remove(r.Context(), "receipt")
	if err := app.renderTemplate(w, r, "virtual-terminal-receipt", &templateData{
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

//login

func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := app.DB.Authenticate(email, password)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
