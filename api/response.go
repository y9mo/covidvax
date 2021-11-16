package api

import (
	"net/http"

	"github.com/y9mo/covidvax/repository"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Msg string `json:"message"`
}

// handleRepositoryError will generate an http status code and a errorResponse
// for a given error from a repository
// it will also log those error if needed
func handleRepositoryError(err error, logger *zap.Logger) (int, *ErrorResponse) {
	if err == repository.ErrRecordNotFound {
		return http.StatusNotFound, &ErrorResponse{Msg: err.Error()}
	}
	if err == repository.ErrUniqueConstraintFailure {
		return http.StatusBadRequest, &ErrorResponse{Msg: "Already exist"}
	}

	if err == repository.ErrInvalidID {
		return http.StatusBadRequest, &ErrorResponse{Msg: err.Error()}
	}

	logger.Error("database Error", zap.Error(err))
	return http.StatusInternalServerError, &ErrorResponse{Msg: "internal error"}
}
