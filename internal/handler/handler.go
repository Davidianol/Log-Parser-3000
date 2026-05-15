package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"log_parser3000/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Parse POST /api/parse/
func (h *Handler) Parse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		writeError(w, http.StatusBadRequest, "body must be JSON with 'path' field")
		return
	}

	logID, err := h.svc.ParseLog(req.Path)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		// битый/некорректный архив
		if logID > 0 {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
				"log_id": logID,
				"error":  err.Error(),
			})
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	start := time.Now()
	elapsed := time.Since(start)

	slog.Info("parse request done",
		"path", req.Path,
		"log_id", logID,
		"duration", elapsed,
	)

	writeJSON(w, http.StatusOK, map[string]any{"log_id": logID})
}

// GetTopology GET /api/topology/{log_id}
func (h *Handler) GetTopology(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	logID, err := extractID(r.URL.Path, "/api/topology/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid log_id")
		return
	}

	topology, err := h.svc.GetTopology(logID)
	if err != nil {
		slog.Error("get topology failed", "log_id", logID, "err", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, topology)
}

// GetNode GET /api/node/{node_id}
func (h *Handler) GetNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nodeID, err := extractID(r.URL.Path, "/api/node/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid node_id")
		return
	}

	node, err := h.svc.GetNode(nodeID)
	if err != nil {
		slog.Error("get node failed", "node_id", nodeID, "err", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, node)
}

// GetPorts GET /api/port/{node_id}
func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nodeID, err := extractID(r.URL.Path, "/api/port/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid node_id")
		return
	}

	ports, err := h.svc.GetPorts(nodeID)
	if err != nil {
		slog.Error("get ports failed", "node_id", nodeID, "err", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ports)
}

// GetLog GET /api/log/{log_id}
func (h *Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	logID, err := extractID(r.URL.Path, "/api/log/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid log_id")
		return
	}

	log, err := h.svc.GetLog(logID)
	if err != nil {
		slog.Error("get log failed", "log_id", logID, "err", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, log)
}

func extractID(path, prefix string) (int, error) {
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.Trim(raw, "/")
	return strconv.Atoi(raw)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response failed", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
