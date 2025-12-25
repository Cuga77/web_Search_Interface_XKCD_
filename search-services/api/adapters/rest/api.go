package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"yadro.com/course/api/adapters/auth"
	"yadro.com/course/api/core"
)

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		replies := make(map[string]string)
		for name, pinger := range pingers {
			if err := pinger.Ping(r.Context()); err != nil {
				replies[name] = "error: " + err.Error()
			} else {
				replies[name] = "ok"
			}
		}

		resp := map[string]interface{}{
			"replies": replies,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Update(r.Context())
		if err != nil {
			if strings.Contains(err.Error(), "update in progress") {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			log.Error("update failed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updater.Stats(r.Context())
		if err != nil {
			log.Error("failed to get stats", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]int{
			"words_total":    stats.WordsTotal,
			"words_unique":   stats.WordsUnique,
			"comics_fetched": stats.ComicsFetched,
			"comics_total":   stats.ComicsTotal,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := updater.Status(r.Context())
		if err != nil {
			log.Error("failed to get status", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"status": string(status),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Drop(r.Context()); err != nil {
			log.Error("failed to drop comics", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewSearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		if phrase == "" {
			http.Error(w, "phrase is required", http.StatusBadRequest)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 10
		if limitStr != "" {
			var err error
			_, err = fmt.Sscanf(limitStr, "%d", &limit)
			if err != nil || limit <= 0 {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}
		}

		result, err := searcher.Search(r.Context(), phrase, limit)
		if err != nil {
			log.Error("failed to search comics", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}

func NewISearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		if phrase == "" {
			http.Error(w, "phrase is required", http.StatusBadRequest)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 10
		if limitStr != "" {
			var err error
			_, err = fmt.Sscanf(limitStr, "%d", &limit)
			if err != nil || limit <= 0 {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}
		}

		result, err := searcher.ISearch(r.Context(), phrase, limit)
		if err != nil {
			log.Error("failed to isearch comics", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}

func NewLoginHandler(log *slog.Logger, auth auth.Authorizer) http.HandlerFunc {
	type loginRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode login request", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		token, err := auth.Login(req.Name, req.Password)
		if err != nil {
			log.Warn("login failed", "user", req.Name, "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(token)); err != nil {
			log.Error("failed to write response", "error", err)
		}
	}
}
