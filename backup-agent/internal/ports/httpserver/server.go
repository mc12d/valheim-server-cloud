package httpserver

import (
	"fmt"
	"net/http"

	"backup-agent/internal/app"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type Server struct {
	BackupJob  app.BackupJob
	RestoreJob app.RestoreJob
	Logger     zerolog.Logger
}

func makeRouter(server Server) chi.Router {
	r := chi.NewRouter()

	r.Post("/backup/make", server.backup)
	r.Post("/backup/restore", server.restore)

	return r
}

const errorBodyFormat = `{"error":"%s"}`

func (s Server) Start(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), makeRouter(s))
}

func (s Server) backup(w http.ResponseWriter, r *http.Request) {
	err := s.BackupJob.Backup()

	if err != nil {
		s.Logger.Error().Err(err).Msgf("httpserver handler on %s", r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(errorBodyFormat, err.Error()))); err != nil {
			s.Logger.Error().Err(err).Msgf("failed to write error response body")
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) restore(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		version = r.URL.Query().Get("version")
	)

	if version != "" {
		err = s.RestoreJob.RestoreVersion(version)
	} else {
		err = s.RestoreJob.Restore()
	}

	if err != nil {
		s.Logger.Error().Err(err).Msgf("httpserver handler on %s", r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(errorBodyFormat, err.Error()))); err != nil {
			s.Logger.Error().Err(err).Msgf("failed to write error response body")
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
