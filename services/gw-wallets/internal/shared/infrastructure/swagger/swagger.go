package swagger

import (
	"net/http"
	"os"

	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
)

type Config struct {
	Title string `mapstructure:"title"`
	Spec  string `mapstructure:"spec"`
}

func Load(cfg Config) (http.HandlerFunc, error) {
	specData, err := os.ReadFile(cfg.Spec)
	if err != nil {
		return nil, err
	}

	return swaggerui.NewHandler(
		swaggerui.WithHTMLTitle(cfg.Title),
		swaggerui.WithSpec(specData),
		swaggerui.WithTryItOutEnabled(true),
	), nil
}

func MustLoad(cfg Config) http.HandlerFunc {
	handler, err := Load(cfg)
	if err != nil {
		panic(err)
	}
	return handler
}
