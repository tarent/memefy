package server

import (
	"encoding/json"
	"fmt"
	"memefy/server/pkg/config"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type AdminHandler struct {
}

type linksInfo struct {
	Links map[string]link `json:"_links"`
}

type link struct {
	Href   string   `json:"href"`
	Method []string `json:"method"`
}

//IndexHandler returns a service self description
func (h *AdminHandler) IndexHandler(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links := make(map[string]link, 20)

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			name := route.GetName()
			href, _ := route.GetPathTemplate()
			methods, _ := route.GetMethods()

			links[name] = link{
				Href:   href,
				Method: methods,
			}

			return nil
		})

		linksInfos := linksInfo{
			Links: links,
		}

		js, err := json.Marshal(linksInfos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *AdminHandler) AdminInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminInfoTemplate := `{
			"build": {
				"BuildNumber": "%s",
				"CommitShortHash": "%s",
				"BuildTime": "%s"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, adminInfoTemplate, config.BuildNumber, config.CommitShortHash, config.BuildTime)
	}
}

// HealthCheckHandler checks if the service can handle traffic and returns 200 or 503
func (h *AdminHandler) HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// -----

func (h *AdminHandler) PostMemeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}