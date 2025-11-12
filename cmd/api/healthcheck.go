package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := map[string]any{
		"status":      "available",
		"environment": app.cfg.env,
		"version":     version,
	}

	err := json.NewEncoder(w).Encode(js)
	if err != nil {
		http.Error(w, "something went wrong with the server", http.StatusInternalServerError)
	}
}
