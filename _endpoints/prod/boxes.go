package handler

import (
	"github.com/gofiber/fiber/v2"
	domainService "gitlab.com/shipink/box-api/boxes/service"
)

const InputMinLength = 1
const InputMaxLength = 255

type Handler struct {
	service domainService.Service
}

func New(service domainService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	// krakend:Role:admin,user
	// krakend:Method:GET
	// krakend:Endpoint:/boxes
	// krakend:ServiceName: box-api
	// krakend:QueryStrings:user_id,page,limit,sort
	// krakend:RateLimit:header,100
	app.Get("/boxes", h.Index)

	// krakend:Role:admin,user
	// krakend:Method:POST
	// krakend:Endpoint:/boxes
	// krakend:ServiceName: box-api
	app.Post("/boxes", h.Create)

	// krakend:Role:admin,user
	// krakend:Method:GET
	// krakend:Endpoint:/boxes/{box_id}
	// krakend:ServiceName: box-api
	// krakend:RateLimit:ip,150
	app.Get("/boxes/:id", h.Get)

	// krakend:Role:admin,user
	// krakend:Method:PUT
	// krakend:Endpoint:/boxes/{box_id}
	// krakend:ServiceName: box-api
	// krakend:RateLimit:header,200
	app.Put("/boxes/:id", h.Update)

	// krakend:Role:admin,user
	// krakend:Method:DELETE
	// krakend:Endpoint:/box/{box_id}
	// krakend:ServiceName: box-api
	// krakend:RateLimit:header,200
	app.Delete("/boxes/:id", h.Delete)
}
