package restapi

import (
	"encoding/json"
	"net/http"
	"protchain/internal/schema"
	"protchain/internal/value"

	"github.com/go-chi/chi/v5"
)

func (a *API) ProteinRoutes(router *chi.Mux) http.Handler {
	router.Method(http.MethodPost, "/", Handler(a.GetProteinH))

	return router
}

func (a *API) GetProteinH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.GetProteinReq

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, nil)
	}

	res, err := a.RetrieveProteinDetail(reqPayload)
	if err != nil {
		return respondWithError(err, "failed to retrieve protein. Please try again", value.Error, nil)
	}

	return &ServerResponse{
		Message:    "protein retrieved",
		Status:     value.Success,
		StatusCode: http.StatusOK,
		Payload:    res,
		File:       res.File,
		FileName:   res.PrimaryAccession,
		FileType:   "chemical/x-pdb",
	}
}
