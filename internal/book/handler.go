package book

import (
	"errors"
	"strconv"

	"restapi/internal/response"

	"github.com/gofiber/fiber/v2"
)

// Handler exposes book endpoints over HTTP.
type Handler struct {
	svc *Service
}

// NewHandler returns a Handler backed by svc.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register wires book routes onto r. Mutating routes are wrapped with auth.
func (h *Handler) Register(r fiber.Router, auth fiber.Handler) {
	r.Get("/book", h.list)
	r.Get("/book/:id", h.get)
	r.Post("/book", auth, h.create)
	r.Patch("/book/:id", auth, h.update)
	r.Delete("/book/:id", auth, h.delete)
}

func (h *Handler) list(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	books, err := h.svc.List(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to list books")
	}
	return c.JSON(books)
}

func (h *Handler) get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	b, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return h.mapError(c, err)
	}
	return c.JSON(b)
}

func (h *Handler) create(c *fiber.Ctx) error {
	var b Book
	if err := c.BodyParser(&b); err != nil {
		return response.Error(c, fiber.StatusNotAcceptable, "invalid request body")
	}
	if err := h.svc.Create(c.Context(), &b); err != nil {
		return h.mapError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(b)
}

func (h *Handler) update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	var input Book
	if err := c.BodyParser(&input); err != nil {
		return response.Error(c, fiber.StatusNotAcceptable, "invalid request body")
	}
	b, err := h.svc.Update(c.Context(), id, input)
	if err != nil {
		return h.mapError(c, err)
	}
	return c.JSON(b)
}

func (h *Handler) delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return h.mapError(c, err)
	}
	return c.JSON(fiber.Map{"message": "book deleted"})
}

// mapError translates domain errors to HTTP responses.
func (h *Handler) mapError(c *fiber.Ctx, err error) error {
	var ve ValidationError
	switch {
	case errors.Is(err, ErrNotFound):
		return response.Error(c, fiber.StatusNotFound, ErrNotFound.Error())
	case errors.As(err, &ve):
		return response.Error(c, fiber.StatusBadRequest, ve.Msg)
	default:
		return response.Error(c, fiber.StatusInternalServerError, "internal server error")
	}
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
