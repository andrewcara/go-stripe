package main

import (
	"net/http"
	"os"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	pk_key := os.Getenv("pk_test")
	data := map[string]interface{}{
		"pk_key": pk_key,
	}
	td := &templateData{Data: data}

	if err := app.renderTemplate(w, r, "terminal", td); err != nil {
		app.errorLog.Println(err)
	}
}
