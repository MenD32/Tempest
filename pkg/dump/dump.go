package dump

import (
	"github.com/MenD32/Tempest/pkg/response"
)

type Dumper interface {
	Dump([]response.Response) error
}
