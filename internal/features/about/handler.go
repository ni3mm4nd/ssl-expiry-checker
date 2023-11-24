package about

import (
	"net/http"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/templates/pages"
)

func About(w http.ResponseWriter, r *http.Request) {
	if err := pages.AboutPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
