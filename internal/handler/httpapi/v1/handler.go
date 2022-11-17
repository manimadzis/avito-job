package v1

import (
	"avito-job/internal/domain"
	"avito-job/internal/repository"
	"avito-job/internal/service"
	"avito-job/pkg/logging"

	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

type Handler struct {
	router  *httprouter.Router
	service service.Service
	logger  logging.Logger
}

func NewHandler(router *httprouter.Router, service service.Service, logger logging.Logger) *Handler {
	h := &Handler{
		router:  router,
		service: service,
		logger:  logger,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	h.router.POST("/v1/user/:user_id/reserve", h.reserveBalance)
	h.router.POST("/v1/user/:user_id/recognize", h.recognizeRevenue)
	h.router.GET("/v1/user/:user_id/balance", h.getBalance)
	h.router.POST("/v1/user/:user_id/balance", h.replenishBalance)
	h.router.GET("/v1/user/:user_id/history/:json", h.getHistory)
	h.router.GET("/v1/report/:year/:month", h.getReport)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) reserveBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("reserveBalance handle request %v", r)
	data, err := h.handleBody(w, r)
	if err != nil {
		return
	}

	dto := domain.ReserveMoneyDTO{}

	dto.UserId, err = h.getUserId(ps)

	if err := h.parseBytes(data, &dto); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	err = h.service.ReserveMoney(r.Context(), &dto)
	if err != nil {
		if err == repository.ErrUnknownUser {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			h.logger.Error("reserveBalance: %v", err)
			return
		} else if err == repository.ErrNotEnoughMoney {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			return
		} else if err == repository.ErrTransactionAlreadyExists {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			return
		}
	}

	h.sendResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("getBalance handle request %v", r)
	var err error
	dto := domain.GetBalanceDTO{}
	dto.UserId, err = h.getUserId(ps)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrInvalidUserId.Error()})
		h.logger.Error(ErrInvalidUserId, " ", err)
		return
	}

	if err = dto.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrInvalidUserId.Error()})
		h.logger.Error(ErrInvalidUserId, " ", err)
		return
	}

	money, err := h.service.GetBalance(r.Context(), &dto)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		if err == repository.ErrUnknownUser {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			return
		}
		h.sendResponse(w, http.StatusInternalServerError, nil)
		return
	}

	h.sendResponse(w, http.StatusOK, struct {
		Balance string `json:"balance"`
	}{
		Balance: money.String(),
	})
}

func (h *Handler) replenishBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("replenishBalance handle request %v", r)
	data, err := h.handleBody(w, r)
	if err != nil {
		return
	}

	dto := domain.ReplenishBalanceDTO{}
	dto.UserId, err = h.getUserId(ps)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	if err := h.parseBytes(data, &dto); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	err = h.service.ReplenishBalance(r.Context(), &dto)
	if err != nil {
		h.logger.Errorf("Failed replenish user balance: %v", err)
		if err == repository.ErrUnknownUser {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			return
		}
		h.sendResponse(w, http.StatusInternalServerError, nil)
		return
	}

	h.sendResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("getHistory handle request %v", r)
	var err error
	data := ps.ByName("json")

	dto := domain.GetHistoryDTO{}
	dto.UserId, err = h.getUserId(ps)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrInvalidUserId.Error()})
		h.logger.Error(ErrInvalidUserId, ":", err)
		return
	}

	if err := h.parseBytes([]byte(data), &dto); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	history, err := h.service.GetHistory(r.Context(), &dto)
	if err != nil {
		if err == repository.ErrUnknownUser {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrUnknownUser.Error()})
			return
		}
		h.sendResponse(w, http.StatusInternalServerError, nil)
		h.logger.Error(err)
	}
	if history == nil {
		history = domain.History{}
	}

	h.sendResponse(w, http.StatusOK, history)
}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("Handle request by getReport: %v", r)
	dto := domain.GetMonthlyReportDTO{}
	var err error
	for _, param := range ps {
		if param.Key == "year" {
			dto.Year, err = strconv.Atoi(param.Value)
			if err != nil {
				h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrYear.Error()})
				h.logger.Errorf("Year convertion failed failed: %v", err)
				return
			}
		}
		if param.Key == "month" {
			dto.Month, err = strconv.Atoi(param.Value)
			if err != nil {
				h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrMonth.Error()})
				h.logger.Errorf("Month convertion failed: %v", err)
				return
			}

		}
	}

	if err := dto.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Errorf("GetMonthlyReportDTO validation failed: %v", err)
		return
	}

	report, err := h.service.GetMonthlyReport(r.Context(), &dto)
	if err != nil {
		h.logger.Errorf("GetMonthlyReport: %v", err)
		h.sendResponse(w, http.StatusInternalServerError, nil)
		return
	}
	if report == nil {
		report = domain.MonthlyReport{}
	}

	h.sendResponse(w, http.StatusOK, struct {
		Report domain.MonthlyReport `json:"report"`
		Length int                  `json:"length"`
	}{
		Report: report,
		Length: len(report),
	})
}

func (h *Handler) recognizeRevenue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("recognizeRevenue handle request %v", r)
	data, err := h.handleBody(w, r)
	if err != nil {
		return
	}

	dto := domain.RecognizeRevenueDTO{}
	dto.UserId, err = h.getUserId(ps)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	if err := h.parseBytes(data, &dto); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	err = h.service.RecognizeRevenue(r.Context(), &dto)
	if err != nil {
		if err == repository.ErrUnknownTransaction {
			h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
			return
		}
		h.logger.Error(err)
		h.sendResponse(w, http.StatusInternalServerError, nil)
	}
	h.sendResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) sendError(w http.ResponseWriter, status int, response ErrorResponse) {
	h.sendResponse(w, status, response)
}

func (h *Handler) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	h.logger.Tracef("sendResponse(%v, %v, %v)", w, status, data)

	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			h.logger.Error("Failed to Marshal errorDTO: %v: data=%v", err, data)
		}
		h.logger.Debugf("json: %s", string(jsonData))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(jsonData)
	}
	w.WriteHeader(status)
}

func (h *Handler) parseBytes(data []byte, dto domain.DTO) error {
	err := json.Unmarshal(data, dto)
	if err != nil {
		h.logger.Errorf("Failed to parse %s: %v", reflect.TypeOf(dto).String(), err)
		return fmt.Errorf("parsing failed: %v", err)
	}
	err = dto.Validate()
	if err != nil {
		h.logger.Errorf("Validation of %s failed: %v", reflect.TypeOf(dto).String(), err)
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}

func (h *Handler) getUserId(ps httprouter.Params) (uint, error) {
	for _, param := range ps {
		if param.Key == "user_id" {
			userId, err := strconv.Atoi(param.Value)
			if err != nil || userId <= 0 {
				return 0, ErrInvalidUserId
			}
			return uint(userId), nil
		}
	}
	return 0, fmt.Errorf("no user_id")
}

func (h *Handler) handleBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(r.Body)
	h.logger.Debugf("Body: %v", string(data))
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, nil)
		h.logger.Errorf("Failed to read body: %v", err)
		return nil, err
	}
	if len(data) == 0 {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrEmptyBody.Error()})
		h.logger.Error(ErrEmptyBody)
		return nil, ErrEmptyBody
	}
	return data, nil
}
