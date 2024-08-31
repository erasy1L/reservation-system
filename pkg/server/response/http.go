package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type BadRequestResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Invalid input provided"`
	Data    any    `json:"data"`
}

type InternalServerErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"An unexpected error occurred"`
	Data    any    `json:"data"`
}

type BaseObject struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
} // @Response

func OK(w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, http.StatusOK)

	v := BaseObject{
		Success: true,
		Data:    data,
	}
	render.JSON(w, r, v)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, data any) {
	render.Status(r, http.StatusBadRequest)

	v := BadRequestResponse{
		Success: false,
		Data:    data,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)

	v := InternalServerErrorResponse{
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Conflict(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
}

func Created(w http.ResponseWriter, r *http.Request, ID string) {
	w.Header().Set("Location", r.URL.Path+"/"+ID)
	w.WriteHeader(http.StatusCreated)
}
