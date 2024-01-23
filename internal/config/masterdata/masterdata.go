package masterdata

import (
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

const defaultPath = "./config/"

func LoadConfig() masterdata.Config {

	return masterdata.Config{
		Path: defaultPath,
	}
}
