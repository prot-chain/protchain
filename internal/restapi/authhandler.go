package restapi

import (
	"encoding/json"
	"net/http"
	"protchain/internal/schema"
	"protchain/internal/value"

	"github.com/go-chi/chi/v5"
)

func (a *API) AuthRoutes() http.Handler {
	router := chi.NewRouter()

	router.Method(http.MethodPost, "/register", Handler(a.RegisterH))
	router.Method(http.MethodPost, "/login", Handler(a.LoginH))
	router.Method(http.MethodPost, "/google-oauth", Handler(a.GoogleOAuthH))

	return router
}

func (a *API) RegisterH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, nil)
	}

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
