package bioapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"protchain/internal/config"
	"protchain/internal/schema"

	"github.com/pkg/errors"
)

type Client struct {
	BaseUrl string
}

func NewClient(config *config.Config) *Client {
	return &Client{
		BaseUrl: config.BioAPIUrl,
	}
}

func (c *Client) RetrieveProtein(payload schema.GetProteinReq) (schema.ProteinData, error) {
	endpoint := fmt.Sprintf("%s/pdb/protein/%s", c.BaseUrl, payload.Code)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return schema.ProteinData{}, errors.Wrap(err, "failed to create request to retreive protein information from bioapi")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return schema.ProteinData{}, errors.Wrap(err, "failed to execute request with default client to bioapi")
	}

	var response schema.ProteinData
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, errors.Wrap(err, "failed to parse response from bioapi")
	}

	return response, nil
}
