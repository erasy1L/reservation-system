package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"room-reservation/internal/domain/reservation"
	"room-reservation/pkg/log"
	"room-reservation/pkg/router"
	"room-reservation/pkg/server/response"

	_ "room-reservation/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ReservationHandler struct {
	reservationRepo reservation.Repository

	HTTP *chi.Mux
}

// @title Room reservation system
// @version 1.0
// @description This is a simple API project
// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi
func NewReservationHandler(repo reservation.Repository) *ReservationHandler {
	h := &ReservationHandler{reservationRepo: repo}

	h.HTTP = router.New()

	h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

	h.HTTP.Route("/api/v1", func(r chi.Router) {
		r.Mount("/reservations", h.routes())
	})

	return h
}

func (h *ReservationHandler) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", h.createReservation)

	r.Route("/{id}", func(r chi.Router) {
		r.Delete("/", h.deleteReservation)
		r.Patch("/", h.updateReservation)
		r.Get("/", h.getReservation)
	})

	r.Get("/room/{roomID}", h.listRoomReservations)

	return r
}

// @Summary Create new reservation
// @Description Create new reservation
// @Tags Reservations
// @Accept json
// @Param reservation body reservation.Request true "Reservation object to be added"
// @Success 201
// @Failure 409 "Overlapping reservation"
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /reservations [post]
func (h *ReservationHandler) createReservation(w http.ResponseWriter, r *http.Request) {
	logger := log.LoggerFromContext(r.Context())

	var req reservation.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Err(err).Caller().Send()
		response.BadRequest(w, r, err, req)
		return
	}

	if err := req.Validate(); err != nil {
		logger.Err(err).Caller().Send()
		response.BadRequest(w, r, err, req)
		return
	}

	data := reservation.Reservation{
		RoomID:    req.RoomID,
		StartTime: req.StartTime.Time,
		EndTime:   req.EndTime.Time,
	}

	ID, err := h.reservationRepo.Create(r.Context(), data)
	if err != nil {
		if errors.Is(err, reservation.ErrorOverlaps) {
			logger.Err(err).Caller().Send()
			response.Conflict(w)
			return
		}

		logger.Err(err).Caller().Send()
		response.InternalServerError(w, r, err)
		return
	}

	response.Created(w, r, ID)
}

// @Summary List reservations for a room
// @Description List reservations for a room
// @Tags Reservations
// @Accept json
// @Produce json
// @Param roomID path string true "Room id"
// @Success 200 {object} response.BaseObject
// @Success 204
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /reservations/room/{roomID} [get]
func (h *ReservationHandler) listRoomReservations(w http.ResponseWriter, r *http.Request) {
	logger := log.LoggerFromContext(r.Context())

	roomID := chi.URLParam(r, "roomID")

	data, err := h.reservationRepo.List(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, reservation.ErrorNotFoundForRoom) {
			logger.Err(err).Caller().Send()
			response.NoContent(w)
			return
		}

		logger.Err(err).Caller().Send()
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, reservation.ToResponseSlice(data))
}

// @Summary Get individual reservation
// @Description Get individual reservation
// @Tags Reservations
// @Accept json
// @Param id path string true "Reservation id"
// @Success 200
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /reservations/{id} [get]
func (h *ReservationHandler) getReservation(w http.ResponseWriter, r *http.Request) {
	logger := log.LoggerFromContext(r.Context())

	ID := chi.URLParam(r, "id")

	data, err := h.reservationRepo.Get(r.Context(), ID)
	if err != nil {
		if errors.Is(err, reservation.ErrorNotFound) {
			logger.Err(err).Caller().Send()
			response.BadRequest(w, r, err, ID)
			return
		}

		logger.Err(err).Caller().Send()
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, data)
}

// @Summary Delete reservation
// @Description Delete reservation
// @Tags Reservations
// @Accept json
// @Param id path string true "Reservation id"
// @Success 204
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /reservations/{id} [delete]
func (h *ReservationHandler) deleteReservation(w http.ResponseWriter, r *http.Request) {
	logger := log.LoggerFromContext(r.Context())

	ID := chi.URLParam(r, "id")

	err := h.reservationRepo.Delete(r.Context(), ID)
	if err != nil {
		if errors.Is(err, reservation.ErrorNotFound) {
			logger.Err(err).Caller().Send()
			response.BadRequest(w, r, err, ID)
			return
		}

		logger.Err(err).Caller().Send()
		response.InternalServerError(w, r, err)
		return
	}

	response.NoContent(w)
}

// @Summary Update reservation
// @Description Update reservation
// @Tags Reservations
// @Accept json
// @Param id path string true "Reservation id"
// @Param body body reservation.UpdateRequest true "Reservation details"
// @Success 204
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /reservations/{id} [patch]
func (h *ReservationHandler) updateReservation(w http.ResponseWriter, r *http.Request) {
	logger := log.LoggerFromContext(r.Context())

	ID := chi.URLParam(r, "id")

	var req reservation.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Err(err).Caller().Send()
		response.BadRequest(w, r, err, req)
		return
	}

	if err := req.Validate(); err != nil {
		logger.Err(err).Caller().Send()
		response.BadRequest(w, r, err, req)
		return
	}

	data := reservation.Reservation{
		RoomID:    req.RoomID,
		StartTime: req.StartTime.Time,
		EndTime:   req.EndTime.Time,
	}

	err := h.reservationRepo.Update(r.Context(), ID, data)
	if err != nil {
		if errors.Is(err, reservation.ErrorNotFound) {
			logger.Err(err).Caller().Send()
			response.BadRequest(w, r, err, ID)
			return
		}

		logger.Err(err).Caller().Send()
		response.InternalServerError(w, r, err)
		return
	}

	response.NoContent(w)
}
