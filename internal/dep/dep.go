package dep

import (
	"protchain/blockchain/gateway"
	"protchain/internal/config"
	"protchain/internal/dal"
	"protchain/internal/service/bioapi"
	"protchain/internal/storage"
	"sync"
)

var (
	dep  *Dependencies
	once sync.Once
)

type Dependencies struct {
	// Services
	Bioapi       *bioapi.Client
	IPFS         *storage.IPFSStorage
	FabricClient *gateway.FabricClient

	// DAL
	DAL *dal.DAL
}

// New initializes the dependencies required for
// the application to function
func New(cfg *config.Config) *Dependencies {
	once.Do(func() {
		dep = &Dependencies{
			Bioapi: bioapi.NewClient(cfg),
			IPFS:   storage.NewIPFSStorage(cfg),
			DAL:    dal.NewDAL(cfg),
		}
	})

	return dep
}
