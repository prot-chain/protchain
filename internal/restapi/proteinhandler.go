package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"protchain/internal/schema"
	"protchain/internal/tracing"
	"protchain/internal/value"

	"github.com/go-chi/chi/v5"
)

func (a *API) ProteinRoutes() http.Handler {
	router := chi.NewRouter()

	router.Method(http.MethodGet, "/retrieve-protein", Handler(a.GetProteinH))

	return router
}

func (a *API) GetProteinH(w http.ResponseWriter, r *http.Request) *ServerResponse {
	var reqPayload schema.GetProteinReq
	tc := r.Context().Value(value.ContextTracingKey).(tracing.Context)

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		return respondWithError(err, "bad request body", value.BadRequest, &tc)
	}

	res, err := a.Deps.Bioapi.RetrieveProtein(reqPayload)
	if err != nil {
		return respondWithError(err, "failed to retrieve protein. Please try again", value.Error, &tc)
	}

	fmt.Println(res)

	return &ServerResponse{
		Message:    "protein retrieved",
		Status:     value.Success,
		StatusCode: http.StatusOK,
	}
}
