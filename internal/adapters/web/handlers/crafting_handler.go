package handlers

import (
	"encoding/json"
	"net/http"

	"palworld-helper/internal/core/domain"
	"palworld-helper/internal/core/ports"
	"palworld-helper/web/templates"
)

type CraftingHandler struct {
	service ports.CraftingService
}

func NewCraftingHandler(service ports.CraftingService) *CraftingHandler {
	return &CraftingHandler{
		service: service,
	}
}

func (h *CraftingHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(templates.MainPageHTML))
}

func (h *CraftingHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	recipes, err := h.service.GetAllRecipes()
	if err != nil {
		http.Error(w, "Failed to get recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *CraftingHandler) CalculateResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.CraftingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	resourceTotals, err := h.service.CalculateResources(req)
	if err != nil {
		http.Error(w, "Failed to calculate resources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resourceTotals)
}
