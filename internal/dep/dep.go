package dep

import (
	"protchain/internal/config"
	"protchain/internal/dal"
	"sync"
)

var (
	dep  *Dependencies
	once sync.Once
)

type Dependencies struct {
	// Services

	// DAL
	DAL *dal.DAL
}

// New initializes the dependencies required for
// the application to function
func New(cfg *config.Config) *Dependencies {
	once.Do(func() {
		dep = &Dependencies{
			DAL: dal.NewDAL(cfg),
		}
	})

	return dep
}
