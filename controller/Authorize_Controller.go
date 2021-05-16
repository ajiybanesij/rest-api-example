package controller

import (
	"encoding/json"
	"net/http"
	"rest-api-example/middleware"
	"rest-api-example/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Login
	_ = json.NewDecoder(r.Body).Decode(&user)
	token, err := middleware.GenerateJWT(user)
	var token_struct models.Token
	if err != nil {
		token_struct.Message = "Error"
		token_struct.Token = ""
	}
	token_struct.Token = token
	token_struct.Message = "Success"

	json.NewEncoder(w).Encode(token_struct)
}
