package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/config"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/features/about"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/features/configview"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/features/home"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/templates/partials"
)

func main() {
	var port = ":3000"
	flag.StringVar(&port, "port", port, "port to listen on")
	flag.Parse()

	// Check if there is any data available from stdin
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", home.Home)
	r.Get("/about", about.About)
	r.Get("/sslchecks/view", partials.RenderSSLChecksController)
	r.Get("/sslchecks/rescan", partials.RescanAll)
	r.Get("/sslchecks/check/{target}", partials.RenderSSLCheckController)
	r.Delete("/sslchecks/{target}", partials.DeleteSSLCheck)
	r.Post("/sslchecks", partials.AddSSLCheck)
	r.Get("/config", configview.Config)
	r.Get("/config/view", partials.GetConfig)

	server := &http.Server{
		Addr:    port,
		Handler: http.TimeoutHandler(r, 30*time.Second, "request timed out"),
	}

	// Display the localhost address and port
	fmt.Printf("Listening on http://localhost%s\n", port)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.Init()
}
