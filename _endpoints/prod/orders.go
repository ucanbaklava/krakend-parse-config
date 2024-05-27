package handler

import (
	"github.com/gofiber/fiber/v2"
	domainService "gitlab.com/shipink/template-api/examples/service"
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
	// krakend:Endpoint:/orders
	// krakend:ServiceName: order-api
	// krakend:QueryStrings:user_id,page,limit,sort
	app.Get("/orders", h.Index)

	// krakend:Role:admin,user
	// krakend:Method:POST
	// krakend:Endpoint:/orders
	// krakend:ServiceName: order-api
	app.Post("/orders", h.Create)

	// krakend:Role:admin,user
	// krakend:Method:GET
	// krakend:Endpoint:/orders/{order_id}
	// krakend:ServiceName: order-api
	app.Get("/orders/:id", h.Get)

	// krakend:Role:admin,user
	// krakend:Method:PUT
	// krakend:Endpoint:/orders/{order_id}
	// krakend:ServiceName: order-api
	app.Put("/orders/:id", h.Update)

	// krakend:Role:admin,user
	// krakend:Method:DELETE
	// krakend:Endpoint:/orders/{order_id}
	// krakend:ServiceName: order-api
	app.Delete("/orders/:id", h.Delete)
}
