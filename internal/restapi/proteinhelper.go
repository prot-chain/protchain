package restapi

import (
	"fmt"
	"io"
	"net/http"
	"protchain/internal/schema"
	"protchain/pkg/function"

	"github.com/cloudflare/cfssl/log"
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

	if bcRes != "" {
		print("data stored on blockchain already -> ", bcRes)
		return schema.Protein{}, nil
	}

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
	fmt.Println(res.Data)

	result, err := a.Deps.FabricClient.SubmitTransaction("StoreMetadata", res.Data.GenerateContractArgs("StoreMetadata")...)
	if err != nil {
		return schema.Protein{}, errors.Wrap(err, "Failed to submit transaction")
	}
	log.Info("blockchain result", result)
	return res.Data, nil
}
