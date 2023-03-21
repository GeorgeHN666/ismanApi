package main

import (
	"fmt"
	"net/http"
)

func (app *application) contactoEP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var res struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	var cont *Contact
	err := app.ReadJSON(w, r, &cont)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	if len(cont.Email) < 1 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	} else if len(cont.Name) < 1 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	} else if len(cont.Desc) < 1 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	} else {
		err = app.SendContactForm(cont.Name, cont.Email, cont.Desc)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	res.Error = false
	res.Message = "all good"
	app.WriteJSON(w, r, res, http.StatusOK)

}

func (app *application) SubscribeUserEP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var res struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	var sub *Subscriber
	err := app.ReadJSON(w, r, &sub)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	us, err := app.GetSubscriber(sub.Email)
	if err == nil {
		// If user its subscribed
		res.Error = true
		res.Message = fmt.Sprintf("Hola %s, Ya estas subscrito!!", us.FullName)
		app.WriteJSON(w, r, res, http.StatusFound)
		return
	}

	err = app.InsertSubscriber(sub)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	err = app.ReportSubscription(sub.Email, sub.FullName)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	res.Error = false
	res.Message = fmt.Sprintf("%s exitosamente subscrito", sub.FullName)

	app.WriteJSON(w, r, res, http.StatusCreated)

}

func (app *application) DeleteSubscriptionEP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var res struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	name := r.URL.Query().Get("au_s")
	email := r.URL.Query().Get("e_so")
	to := r.URL.Query().Get("TMS_A")

	us, err := app.GetSubscriber(email)
	if err != nil {
		// If user its subscribed
		res.Error = true
		res.Message = fmt.Sprintf("El usuario %s que intenta buscar ya esta dado de baja", us.FullName)
		app.WriteJSON(w, r, res, http.StatusFound)
		return
	}

	if to == us.Code {
		err = app.DeleteSubscription(us)
		if err != nil {
			app.BadRequest(w, r, err)
			return
		}
	} else {
		res.Error = true
		res.Message = fmt.Sprintf("El usuario %s no ha sido encontrado, verifica que el link sea correcto o pide uno nuevo", name)
		app.WriteJSON(w, r, res, http.StatusFound)
	}

	res.Error = false
	res.Message = fmt.Sprintf("Usuario %s exitosamente dado de baja", us.FullName)

	app.WriteJSON(w, r, res, http.StatusOK)

}
