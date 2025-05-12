package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andrewcara/go-stripe.git/internal/cards"
	"github.com/andrewcara/go-stripe.git/internal/models"
	"github.com/go-chi/chi/v5"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Content string `json:"content"`
	ID      int    `json:"id"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {

	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)

	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", " ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {

		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", " ")

		if err != nil {
			app.errorLog.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

}

func (app *application) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	eventID, _ := strconv.Atoi(id)

	event, err := app.DB.GetTicketEvents(eventID)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(event, "", " ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &userInput)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	validPassword, err := app.passwordMatches(user.Password, userInput.Password)

	if err != nil {
		app.invalidCredentials(w)
		return
	}

	if !validPassword {
		app.invalidCredentials(w)
		return
	}

	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)

	if err != nil {
		app.badRequest(w, r, err)
	}

	//save to db

	err = app.DB.InsertToken(token, user)

	if err != nil {
		app.badRequest(w, r, err)
	}

	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}

	payload.Error = false
	payload.Message = "Success"
	payload.Token = token

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) AuthenticateToken(r *http.Request) (*models.User, error) {
	var u models.User

	return &u, nil
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {

	//validate token and get user

	user, err := app.AuthenticateToken(r)

	if err != nil {
		app.invalidCredentials(w)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"string"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated yser %s", user.Email)
	app.writeJSON(w, http.StatusOK, payload)

}
