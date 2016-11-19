package webhooks

import (
	"net/http"

	"github.com/inconshreveable/log15"
	hobware "github.com/leeola/hob/middleware"
	"github.com/pressly/chi"
)

func ListenAndServe(bind string) {
	r := chi.NewRouter()

	// Using hobs middleware for convenience.
	r.Use(hobware.Logging(log15.New()))

	r.Post("/webhook", handle())

	http.ListenAndServe(bind, r)
}
