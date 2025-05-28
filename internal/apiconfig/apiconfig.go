package apiconfig

import (
	"sync/atomic"

	"github.com/bbarrington0099/Chirpy/internal/database"
)

type Conf struct {
	Port string
	FilepathRoot string
	FileserverHits atomic.Int32
	QueryCollection *database.Queries
	Platform string
}