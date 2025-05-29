package api

import (
	"encoding/json"
	"net/http"

	"github.com/liangsj/vimcoplit/internal/core"
)

// ContextHandler 提供上下文管理的 HTTP API
func ContextHandler(svc core.Service) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/context/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			ID    string           `json:"id"`
			Type  core.ContextType `json:"type"`
			Value string           `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		item := core.NewContextItem(req.ID, req.Type, req.Value)
		svc.GetContextManager().AddItem(item)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/context/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := svc.GetContextManager().RemoveItem(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/context/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		item, err := svc.GetContextManager().GetItem(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(item)
	})

	mux.HandleFunc("/api/context/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		items := svc.GetContextManager().ListItems()
		json.NewEncoder(w).Encode(items)
	})

	return mux
}
