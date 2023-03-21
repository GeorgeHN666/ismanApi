package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	maxBytes := int64(1048576)

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("something went wrong")
	}
	return nil
}

func (app *application) WriteJSON(w http.ResponseWriter, r *http.Request, data interface{}, status int, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
	return nil

}

func (app *application) BadRequest(w http.ResponseWriter, r *http.Request, err error) error {

	var response struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	response.Error = true
	response.Message = err.Error()

	out, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(out)
	return nil
}

func (app *application) Hash(value string) (string, error) {
	cost := 8

	hash, err := bcrypt.GenerateFromPassword([]byte(value), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (app *application) DeHash(hash string, value string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// GenerateCode Will generate a code to a given user
func (app *application) GenerateCode() (string, error) {
	length := 25

	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}
