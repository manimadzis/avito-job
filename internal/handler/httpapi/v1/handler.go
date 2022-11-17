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
	h.router.GET("/v1/user/:user_id/balance", h.getBalance)
	h.router.POST("/v1/user/:user_id/balance", h.replenishBalance)
	h.router.GET("/v1/user/:user_id/history", h.getHistory)
	h.router.GET("/v1/report/:year/:month", h.getReport)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) reserveBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("getBalance handle request %v", r)
	data, err := io.ReadAll(r.Body)
	h.logger.Debugf("Body: %v", string(data))

	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, nil)
		h.logger.Errorf("Failed to read body: %v", err)
		return
	}

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
			h.sendError(w, http.StatusNotFound, ErrorResponse{Msg: err.Error()})
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
	data, err := io.ReadAll(r.Body)
	h.logger.Debugf("Body: %v", string(data))
	h.logger.Debugf("Params: %v", ps)
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
			h.sendError(w, http.StatusNotFound, ErrorResponse{Msg: err.Error()})
			return
		}
		h.sendResponse(w, http.StatusInternalServerError, nil)
		return
	}

	h.sendResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Tracef("getHistory handle request %v", r)
	data, err := io.ReadAll(r.Body)
	h.logger.Debugf("Body: %v", string(data))

	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, nil)
		h.logger.Errorf("Failed to read body: %v", err)
		return
	}

	dto := domain.GetHistoryDTO{}
	dto.UserId, err = h.getUserId(ps)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: ErrInvalidUserId.Error()})
		h.logger.Error(ErrInvalidUserId, ":", err)
		return
	}

	if err := h.parseBytes(data, &dto); err != nil {
		h.sendError(w, http.StatusBadRequest, ErrorResponse{Msg: err.Error()})
		h.logger.Error(err)
		return
	}

	//history, err := h.service.GetHistory(r.Context(), &dto)
	//if err != nil {
	//	h.sendResponse(w, http.StatusInternalServerError, nil)
	//	h.logger.Error(ErrInvalidUserId, " ", err)
	//}

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

func (h *Handler) sendError(w http.ResponseWriter, status int, response ErrorResponse) {
	h.sendResponse(w, status, response)
}

func (h *Handler) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	h.logger.Tracef("sendResponse(%v, %v, %v)", w, status, data)
	w.WriteHeader(status)

	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			h.logger.Error("Failed to Marshal errorDTO: %v: data=%v", err, data)
		}
		h.logger.Debug(string(jsonData))
		w.Write(jsonData)
	}
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
