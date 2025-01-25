package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"protchain/internal/schema"
	"protchain/internal/value"

	"github.com/go-chi/chi/v5"
)

func (a *API) AuthRoutes(r *chi.Mux) http.Handler {
	r.Method(http.MethodPost, "/register", Handler(a.RegisterH))
	r.Method(http.MethodPost, "/login", Handler(a.LoginH))
	r.Method(http.MethodPost, "/google-oauth", Handler(a.GoogleOAuthH))

	return r
}

func (a *API) RegisterH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, nil)
	}
	fmt.Println("payload is -> ", reqPayload)

	res, err := a.RegisterUser(reqPayload)
	if err != nil {
		return respondWithError(err, "registration failed", value.Error, nil)
	}

	return &ServerResponse{
		Message:    "registration successful",
		Status:     value.Success,
		StatusCode: http.StatusOK,
		Payload:    res,
	}
}

func (a *API) LoginH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.LoginReq

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, nil)
	}

	res, err := a.LoginUser(reqPayload)
	if err != nil {
		return respondWithError(err, "login failed", value.Error, nil)
	}

	return &ServerResponse{
		Message:    "login successful",
		Status:     value.Success,
		StatusCode: http.StatusOK,
		Payload:    res,
	}
}

func (a *API) GoogleOAuthH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.GoogleOAuthReq

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, nil)
	}

	res, err := a.GoogleOAuth(reqPayload)
	if err != nil {
		return respondWithError(err, "google oauth failed", value.Error, nil)
	}

	return &ServerResponse{
		Message:    "google oauth successful",
		Status:     value.Success,
		StatusCode: http.StatusOK,
		Payload:    res,
	}
}
