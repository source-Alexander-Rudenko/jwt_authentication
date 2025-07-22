package delivery

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"

	"jwt_auth_project/internal/domain"
	"jwt_auth_project/internal/repo"
	"jwt_auth_project/internal/usecase"
	"jwt_auth_project/internal/utils"
)

// AdsHandler обрабатывает HTTP-запросы для CRUD объявлений
type AdsHandler struct {
	adsUC usecase.AdsUseCase
}

// NewAdsHandler создаёт новый обработчик объявлений
func NewAdsHandler(adsUC usecase.AdsUseCase) *AdsHandler {
	return &AdsHandler{adsUC: adsUC}
}

// RegisterRoutes регистрирует маршруты для работы с объявлениями
func (h *AdsHandler) RegisterRoutes(r *mux.Router) {
	sub := r.PathPrefix("/ads").Subrouter()
	sub.HandleFunc("", h.handleListAds).Methods(http.MethodGet)
	sub.HandleFunc("", h.handleCreateAd).Methods(http.MethodPost)
	sub.HandleFunc("/{id}", h.handleGetAd).Methods(http.MethodGet)
	sub.HandleFunc("/{id}", h.handleUpdateAd).Methods(http.MethodPut)
	sub.HandleFunc("/{id}", h.handleDeleteAd).Methods(http.MethodDelete)
}

// handleCreateAd создаёт новое объявление через multipart/form-data
func (h *AdsHandler) handleCreateAd(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера тела до 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("create ad: parse form failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
		return
	}

	// Чтение и валидация полей
	authorIDStr := r.FormValue("author_id")
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	priceStr := r.FormValue("price")

	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		slog.Error("create ad: invalid author_id", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid author_id"))
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		slog.Error("create ad: invalid price", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid price"))
		return
	}

	// Извлечение файла изображения
	file, header, err := r.FormFile("image")
	if err != nil {
		slog.Error("create ad: image required", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image is required"))
		return
	}
	defer file.Close()

	// Проверка типа Reader
	rdr, ok := file.(io.ReadSeeker)
	if !ok {
		slog.Error("create ad: file not ReadSeeker")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server error"))
		return
	}

	// Формирование payload
	payload := domain.CreateAdPayload{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Price:       price,
		Image:       rdr,
		ImageSize:   header.Size,
		ImageName:   header.Filename,
		ContentType: header.Header.Get("Content-Type"),
	}

	// Вызов бизнес-логики
	ad, err := h.adsUC.CreateAd(r.Context(), payload)
	if err != nil {
		slog.Error("create ad: usecase error", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	slog.Info("ad created", "id", ad.ID)
	utils.WriteJSON(w, http.StatusCreated, ad)
}

// handleGetAd возвращает объявление по UUID
func (h *AdsHandler) handleGetAd(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("get ad: invalid id", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	ad, err := h.adsUC.GetAdByID(r.Context(), id)
	if err != nil {
		slog.Error("get ad: usecase error", "error", err)
		if errors.Is(err, repo.ErrAdNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, ad)
}

// handleListAds возвращает список объявлений с пагинацией и сортировкой
func (h *AdsHandler) handleListAds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := utils.ParseInt(q.Get("limit"), 10)
	offset := utils.ParseInt(q.Get("offset"), 0)
	sortField := q.Get("sort_field")
	sortAsc := q.Get("sort_asc") == "true"

	opts := domain.AdListOptions{
		Limit:     limit,
		Offset:    offset,
		SortField: sortField,
		SortAsc:   sortAsc,
	}

	ads, err := h.adsUC.ListAds(r.Context(), opts)
	if err != nil {
		slog.Error("list ads: usecase error", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ads)
}

// handleUpdateAd обновляет объявление и опционально заменяет картинку
func (h *AdsHandler) handleUpdateAd(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("update ad: invalid id", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	// Ограничение размера тела
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("update ad: parse form failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
		return
	}

	// Чтение полей
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		slog.Error("update ad: invalid price", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid price"))
		return
	}

	// Чтение файла (необязательно)
	var rdr io.ReadSeeker
	var header *multipart.FileHeader
	if file, fh, err := r.FormFile("image"); err == nil {
		defer file.Close()
		rs, ok := file.(io.ReadSeeker)
		if !ok {
			slog.Error("update ad: file not ReadSeeker")
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server error"))
			return
		}

		rdr = rs
		header = fh
	}

	// Подготовка payload
	payload := domain.UpdateAdPayload{
		ID:          id,
		Title:       title,
		Description: description,
		Price:       price,
		Image:       rdr,
	}
	if rdr != nil {
		payload.ImageSize = header.Size
		payload.ImageName = header.Filename
		payload.ContentType = header.Header.Get("Content-Type")
	}

	ad, err := h.adsUC.UpdateAd(r.Context(), payload)
	if err != nil {
		slog.Error("update ad: usecase error", "error", err)
		if errors.Is(err, repo.ErrAdNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusBadRequest, err)
		}
		return
	}

	slog.Info("ad updated", "id", ad.ID)
	utils.WriteJSON(w, http.StatusOK, ad)
}

// handleDeleteAd удаляет объявление по UUID
func (h *AdsHandler) handleDeleteAd(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("delete ad: invalid id", "error", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	err = h.adsUC.DeleteAd(r.Context(), id)
	if err != nil {
		slog.Error("delete ad: usecase error", "error", err)
		if errors.Is(err, repo.ErrAdNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	slog.Info("ad deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
