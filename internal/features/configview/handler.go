package configview

import (
	"net/http"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/templates/pages"
)

func Config(w http.ResponseWriter, r *http.Request) {
	if err := pages.ConfigPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
