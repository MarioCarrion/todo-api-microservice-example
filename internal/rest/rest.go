package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.opentelemetry.io/otel"

	"github.com/MarioCarrion/todo-api/internal"
)

const otelName = "github.com/MarioCarrion/todo-api/internal/rest"

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
}

func renderErrorResponse(router *gin.Context, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var ierr *internal.Error
	if !errors.As(err, &ierr) {
		resp.Error = "internal error"
	} else {
		switch ierr.Code() {
		case internal.ErrorCodeNotFound:
			status = http.StatusNotFound
		case internal.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest

			var verrors validation.Errors
			if errors.As(ierr, &verrors) {
				resp.Validations = verrors
			}
		case internal.ErrorCodeUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
		}
	}

	if err != nil {
		_, span := otel.Tracer(otelName).Start(router, "renderErrorResponse")
		defer span.End()

		span.RecordError(err)
	}

	// XXX fmt.Printf("Error: %v\n", err)
	router.JSON(status, resp)
}
