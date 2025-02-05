package errors

import (
	"fmt"
	"runtime/debug"

	"database/sql"
	"errors"
	"net/http"

	"github.com/courage173/go-auth-api/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)


func Handler(logger log.Logger) routing.Handler{
	return func(c *routing.Context) (err error){
		defer func(){
			l := logger.With(c.Request.Context())

			if e := recover(); e != nil {
				var ok bool

				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}

				l.Errorf("recovered from panic (%v): %s", err, debug.Stack())

			}

			if err != nil {
				res := buildErrorResponse(err)
				if res.StatusCode() == http.StatusInternalServerError {
					l.Errorf("encountered internal server error: %v", err)
				}
				c.Response.WriteHeader(res.StatusCode())
				if err = c.Write(res); err != nil {
					l.Errorf("failed writing error response: %v", err)
				}
				c.Abort() 
				err = nil 
			}
		}()

		return c.Next()
	}
}

func buildErrorResponse(err error) ErrorResponse {
	switch err.(type){
		case ErrorResponse:
			return err.(ErrorResponse)
        case validation.Errors:
			return InvalidInput(err.(validation.Errors))
		case routing.HTTPError:
			switch err.(routing.HTTPError).StatusCode() {
				case http.StatusNotFound:
                    return NotFound("")
                default:
					return ErrorResponse{
						Status:  err.(routing.HTTPError).StatusCode(),
						Message: err.Error(),
					}
			}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NotFound("")
	}
	return InternalServerError("")
}