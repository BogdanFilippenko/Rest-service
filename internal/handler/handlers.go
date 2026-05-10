package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"service/internal/model"
	"service/internal/repository"
	"strconv"
	"github.com/google/uuid"
	
	"service/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *repository.Repository
	svc  *service.Service
	log  *slog.Logger
}

func New(repo *repository.Repository, svc *service.Service, log *slog.Logger) *Handler {
	return &Handler{repo: repo, svc: svc, log: log}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions", h.ListSubscriptions)
	r.Get("/subscriptions/{id}", h.GetSubscription)
	r.Put("/subscriptions/{id}", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)
	r.Post("/subscriptions/cost", h.CalculateCost)

	return r
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	
	if sub.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}
	if sub.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(r.Context(), sub)
	if err != nil {
		h.log.Error("failed to create", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.repo.List(r.Context())
	if err != nil {
		h.log.Error("failed to list", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(subs)
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	sub, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	var sub model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if sub.Price <= 0 || sub.ServiceName == "" {
		http.Error(w, "invalid input data", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), id, sub); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type CostRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

func (h *Handler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	var req CostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	cost, err := h.svc.CalculateCost(r.Context(), req.UserID, req.ServiceName, req.StartDate, req.EndDate)
	if err != nil {
		h.log.Error("failed to calculate cost", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total_cost": cost})
}
