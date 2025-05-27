package apiConfig

import (
	"sync/atomic"
)

type Conf struct {
	Port string
	FilepathRoot string
	FileserverHits atomic.Int32
}