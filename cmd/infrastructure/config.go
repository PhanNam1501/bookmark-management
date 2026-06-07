package infrastructure

import "github.com/PhanNam1501/bookmark-management/internal/api"

// InitConfig loads application configuration
func InitConfig() (*api.Config, error) {
	cfg, err := api.NewConfig("")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
