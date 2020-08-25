package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"yal/models"
	u "yal/utils"

	"github.com/gorilla/mux"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := new(models.Account)
	// decode body to struct, exit on error
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		u.Respond(w, u.Message(false, "Invalid request: "+err.Error()))
		return
	}

	// create acc
	resp := account.Create()
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := new(models.Account)

	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		u.Respond(w, u.Message(false, "Invalid request: "+err.Error()))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	// we've got id from request string
	resp := models.GetUser(uint(id))

	u.Respond(w, resp)
}

var GetUsers = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetUsers()
	u.Respond(w, resp)
}

var UpdateAccount = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	// get update params from body
	updateData := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		u.Respond(w, u.Message(false, "Invalid request: "+err.Error()))
		return
	}
	log.Printf("We've got params to update: %+v", updateData)

	var resp map[string]interface{}
	// we've got id from request string
	resp = models.GetUser(uint(id))

	for k, v := range updateData {
		switch k {
		case "email":
			resp["account"].(*models.Account).Email = v.(string)
		case "password":
			resp["account"].(*models.Account).Password = v.(string)
		case "about":
			resp["account"].(*models.Account).About = v.(string)
		default:
			log.Printf("Field %v is not available to update", k)
		}
	}

	// get updated object and process update
	account := resp["account"].(*models.Account)
	if resp := account.Update(); resp == nil {
		u.Respond(w, u.Message(false, "No Update response"))
		return
	}

	// return updated response
	resp = models.GetUser(uint(id))

	u.Respond(w, resp)
}
