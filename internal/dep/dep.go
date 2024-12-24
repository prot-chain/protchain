package dep

import (
	"protchain/internal/config"
	"protchain/internal/dal"
	"protchain/internal/service/bioapi"
	"sync"
)

var (
	dep  *Dependencies
	once sync.Once
)

type Dependencies struct {
	// Services
	Bioapi *bioapi.Client

	// DAL
	DAL *dal.DAL
}

// New initializes the dependencies required for
// the application to function
func New(cfg *config.Config) *Dependencies {
	once.Do(func() {
		dep = &Dependencies{
			Bioapi: bioapi.NewClient(cfg),
			DAL:    dal.NewDAL(cfg),
		}
	})

	return dep
}
