package partials

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/repository"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service"
)

func RescanAll(w http.ResponseWriter, r *http.Request) {
	service.CheckTargets(service.GetAllTargets())
	RenderSSLChecksController(w, r)
}

func RenderSSLChecksController(w http.ResponseWriter, r *http.Request) {
	checks, err := repository.GetRepos().SSLCheckRepo.ReadAll()
	if err != nil {
		if err := RenderSSLChecks([]sslcheck.SSLCheck{}, err.Error()).Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	if err := RenderSSLChecks(checks, "").Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func DeleteSSLCheck(w http.ResponseWriter, r *http.Request) {
	target := chi.URLParamFromCtx(r.Context(), "target")
	repository.GetRepos().SSLCheckRepo.Delete(target)
	checks, err := repository.GetRepos().SSLCheckRepo.ReadAll()
	if err != nil {
		if err := RenderSSLChecks([]sslcheck.SSLCheck{}, err.Error()).Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	if err := RenderSSLChecks(checks, "").Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func RenderSSLCheckController(w http.ResponseWriter, r *http.Request) {
	target := chi.URLParamFromCtx(r.Context(), "target")
	service.CheckTarget(target)
	check, err := repository.GetRepos().SSLCheckRepo.Read(target)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if err := RenderSSLCheckTable(check).Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func AddSSLCheck(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	target := r.Form.Get("url")
	checks, _ := repository.GetRepos().SSLCheckRepo.ReadAll()
	target, err := formatURL(target)
	if err != nil {
		if err := RenderSSLChecks(checks, "Can not format URL").Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	if target == "" {
		if err := RenderSSLChecks(checks, "target can not be empty").Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	service.CheckTarget(target)
	checks, errChecks := repository.GetRepos().SSLCheckRepo.ReadAll()
	if errChecks != nil {
		if err := RenderSSLChecks([]sslcheck.SSLCheck{}, errChecks.Error()).Render(r.Context(), w); err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}
	if err := RenderSSLChecks(checks, "").Render(r.Context(), w); err != nil {
		w.Write([]byte(err.Error()))
	}

}

func formatURL(targetURL string) (string, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Println(parsedURL)
		return "", errors.New("cannot parse URL")
	}
	if parsedURL.Scheme == "https" {
		targetURL = parsedURL.Host
	}
	if parsedURL.Scheme == "" {
		targetURL = parsedURL.Path
	}
	return targetURL, nil
}
