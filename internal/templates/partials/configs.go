package partials

import (
	"net/http"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/config"
)

func GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := config.Get().GetFullYamlConfig()
	if err != nil {
		if err := ShowConfig("", err.Error()).Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	if err := ShowConfig(string(config), "").Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}
}
