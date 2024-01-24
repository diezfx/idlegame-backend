package supabase

import (
	"github.com/diezfx/idlegame-backend/pkg/auth"
	"github.com/diezfx/idlegame-backend/pkg/configloader"
)

const defaultNamespace = "supabase"

func LoadSupabaseConfig(loader *configloader.Loader) (auth.Config, error) {
	cfg := auth.Config{}

	return cfg, nil
}
