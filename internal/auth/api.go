package auth

import (
	"github.com/courage173/go-auth-api/internal/errors"

	"github.com/courage173/go-auth-api/internal/models"
	"github.com/courage173/go-auth-api/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)


func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.Logger){
	rg.Post("/register", register(service, logger))
	rg.Post("/login", login(service, logger))
}

func register(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error{
		var req models.User

		if err := c.Read(&req); err!= nil {
            logger.With(c.Request.Context()).Errorf("Invalid reqest: %v", err)
            return errors.BadRequest("")
        }

		if err := req.Validate(); err != nil {
			logger.With(c.Request.Context()).Errorf("Invalid reqest: %v", err)
            return err
		}

		response, err := service.Register(c.Request.Context(), req)
		if err!= nil {
            return err
        }

		return c.Write(struct {
			Message string `json:"response"`
		}{response})
	}
}

func login(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error{
		var req models.LoginRequest
    

		
        if err := c.Read(&req); err!= nil {
            logger.With(c.Request.Context()).Errorf("Invalid reqest: %v", err)
            return errors.BadRequest("")
        }

		if err := req.Validate(); err != nil {
			
			logger.With(c.Request.Context()).Errorf("Invalid reqest: %v", err)
            return err
		}

        token, err := service.Login(c.Request.Context(), req.Email, req.Password)
        if err!= nil {
            return err
        }

        return c.Write(struct {
            Token string `json:"token"`
        }{token})
    }
}