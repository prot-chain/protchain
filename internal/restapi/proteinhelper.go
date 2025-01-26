package restapi

import (
	"fmt"
	"io"
	"net/http"
	"protchain/internal/schema"
	"protchain/pkg/function"

	"github.com/pkg/errors"
)

func (a *API) RetrieveProteinDetail(req schema.GetProteinReq) (schema.Protein, error) {
	res, err := a.Deps.Bioapi.RetrieveProtein(req)
	if err != nil {
		return schema.Protein{}, err
	}

	// retrieve PDB file.
	// Check if the file already exists in the blockchain
	bcRes, err := a.Deps.FabricClient.EvaluateTransaction("QueryMetadata", req.Code)
	if err != nil {
		return schema.Protein{}, fmt.Errorf("failed to retrieve PDB file from %s: %w", res.Data.PDBLink, err)
	}

	var payload schema.Protein
	switch bcRes == "" {
	case true: // the protein did not already have a record in the blockchain

		resp, err := http.Get(res.Data.PDBLink)
		if err != nil {
			return schema.Protein{}, fmt.Errorf("failed to retrieve PDB file from %s: %w", res.Data.PDBLink, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return schema.Protein{}, fmt.Errorf("failed to fetch PDB file: received status %d", resp.StatusCode)
		}

		fileData, err := io.ReadAll(resp.Body)
		if err != nil {
			return schema.Protein{}, fmt.Errorf("failed to read PDB file data: %w", err)
		}
		res.Data.FileHash = function.ComputeHash(fileData)

		cid, err := a.Deps.IPFS.Upload(res.Data.PrimaryAccession+".pdb", fileData)
		if err != nil {
			return schema.Protein{}, fmt.Errorf("failed to upload PDB file to IPFS: %w", err)
		}
		res.Data.IPFSCid = cid

		_, err = a.Deps.FabricClient.SubmitTransaction("StoreMetadata", res.Data.GenerateContractArgs("StoreMetadata")...)
		if err != nil {
			return schema.Protein{}, errors.Wrap(err, "Failed to submit transaction")
		}
		payload = res.Data

	default:

		var data map[string]string
		if err := function.Load(bcRes, &data); err != nil {
			return schema.Protein{}, errors.Wrap(err, "Failed to load metadata from blockckain")
		}
		payload.FileHash = data["protein_hash"]
		payload.PrimaryAccession = data["protein_id"]
		payload.IPFSCid = data["file_Url"]
	}

	// retrieve file from IPFS
	rc, err := a.Deps.IPFS.Download(payload.IPFSCid)
	if err != nil {
		return schema.Protein{}, errors.Wrap(err, "Failed to retrieve content from IPFS")
	}
	payload.File = rc

	return payload, nil
}
