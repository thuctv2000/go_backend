package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"my_backend/internal/domain"
)

type LixiHandler struct {
	lixiService domain.LixiService
}

func NewLixiHandler(lixiService domain.LixiService) *LixiHandler {
	return &LixiHandler{
		lixiService: lixiService,
	}
}

// GetActive returns the active lixi config (public endpoint)
func (h *LixiHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	config, err := h.lixiService.GetActiveConfig(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// GetAll returns all lixi configs (admin endpoint)
func (h *LixiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	configs, err := h.lixiService.GetAllConfigs(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return empty array instead of null
	if configs == nil {
		configs = []*domain.LixiConfig{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}

type createLixiRequest struct {
	Name      string                 `json:"name"`
	Envelopes []domain.LixiEnvelope `json:"envelopes"`
}

// Create creates a new lixi config (admin endpoint)
func (h *LixiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createLixiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	config, err := h.lixiService.CreateConfig(r.Context(), req.Name, req.Envelopes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

type updateLixiRequest struct {
	Name      string                 `json:"name"`
	Envelopes []domain.LixiEnvelope `json:"envelopes"`
}

// Update updates a lixi config (admin endpoint)
func (h *LixiHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/admin/lixi/{id}
	id := extractIDFromPath(r.URL.Path, "/api/admin/lixi/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	var req updateLixiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	config, err := h.lixiService.UpdateConfig(r.Context(), id, req.Name, req.Envelopes)
	if err != nil {
		if err.Error() == "lixi config not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// Delete deletes a lixi config (admin endpoint)
func (h *LixiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/admin/lixi/{id}
	id := extractIDFromPath(r.URL.Path, "/api/admin/lixi/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	err := h.lixiService.DeleteConfig(r.Context(), id)
	if err != nil {
		if err.Error() == "lixi config not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "cannot delete active config" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Activate sets a config as active (admin endpoint)
func (h *LixiHandler) Activate(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/admin/lixi/{id}/activate
	path := r.URL.Path
	path = strings.TrimSuffix(path, "/activate")
	id := extractIDFromPath(path, "/api/admin/lixi/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	err := h.lixiService.SetActiveConfig(r.Context(), id)
	if err != nil {
		if err.Error() == "lixi config not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Config activated successfully"})
}

// Helper function to extract ID from path
func extractIDFromPath(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	id := strings.TrimPrefix(path, prefix)
	// Remove trailing slash if present
	id = strings.TrimSuffix(id, "/")
	// Handle paths like "123/activate" - get only the ID part
	if idx := strings.Index(id, "/"); idx != -1 {
		id = id[:idx]
	}
	return id
}
