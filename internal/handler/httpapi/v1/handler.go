package v1

import (
	"avito-job/internal/service"
	"avito-job/pkg/logging"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	h.router.GET("/v1/user/balance/:user_id", h.getBalance)
	h.router.POST("/v1/user/balance/:user_id", h.replenishBalance)
	h.router.GET("/v1/user/history/:user_id", h.getHistory)
	h.router.GET("/v1/report/:year/:month", h.getReport)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) reserveBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) replenishBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *Handler) sendError(w http.ResponseWriter, status int, data interface{}) {
	//h.sendResponse(w, status, errorDTO{error: data})
}

func (h *Handler) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)

	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to Marshal errorDTO: %v: data=%v", err, data)
	}

	w.Write(jsonData)
}
