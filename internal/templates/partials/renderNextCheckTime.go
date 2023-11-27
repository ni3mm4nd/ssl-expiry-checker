package partials

import (
	"net/http"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service"
)

func ShowNextCheckTime(w http.ResponseWriter, r *http.Request) {
	scheduler := service.GetScheduler()
	if err := RenderNextCheckTimePartial(scheduler.Jobs()[0].NextRun().String()).Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}
}
