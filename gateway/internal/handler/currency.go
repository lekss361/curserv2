package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/lekss361/curserv2/gateway/internal/dto"
	errpkg "github.com/lekss361/curserv2/gateway/internal/errors"
	"github.com/lekss361/curserv2/gateway/internal/service"
	"net/http"
	"time"
)

type CurrencyHandler struct {
	currencyService service.CurrencyService
}

func NewCurrencyHandler(cs service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{currencyService: cs}
}

// Expects date in YYYY-MM-DD format.
func (h *CurrencyHandler) Get(w http.ResponseWriter, r *http.Request) {
	dateStr := chi.URLParam(r, "date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	rates, err := h.currencyService.GetRatesByDate(r.Context(), date)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resp := dto.RateResponse{Date: date, Rub: rates}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CurrencyHandler) History(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startStr := q.Get("start")
	endStr := q.Get("end")
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid start date format")
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid end date format")
		return
	}

	history, err := h.currencyService.GetRatesHistory(r.Context(), start, end)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resp := dto.HistoryResponse{History: make([]dto.RateResponse, len(history))}
	for i, dr := range history {
		resp.History[i] = dto.RateResponse{Date: dr.Date, Rub: dr.Rub}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CurrencyHandler) errorResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errpkg.NotFoundError:
		WriteError(w, e.StatusCode(), e.Error())
	default:
		WriteError(w, http.StatusInternalServerError, e.Error())
	}
}
