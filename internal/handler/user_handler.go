package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/yourusername/user-api/internal/models"
	"github.com/yourusername/user-api/internal/service"
	appvalidator "github.com/yourusername/user-api/internal/validator"
)

// UserHandler groups all HTTP handlers for the user resource.
type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
	log      *zap.Logger
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(svc service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: appvalidator.New(),
		log:      log,
	}
}

// -----------------------------------------------------------------------
// POST /users
// -----------------------------------------------------------------------

// CreateUser godoc
// @Summary Create a new user
// @Param body body models.CreateUserRequest true "User payload"
// @Success 201 {object} models.UserResponse
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if errs := h.validate.Struct(req); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{
			Error:   "validation failed",
			Details: formatValidationErrors(errs),
		})
	}

	resp, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		h.log.Error("create user failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// -----------------------------------------------------------------------
// GET /users/:id
// -----------------------------------------------------------------------

// GetUser godoc
// @Summary Get a user by ID
// @Param id path int true "User ID"
// @Success 200 {object} models.UserWithAgeResponse
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	resp, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.log.Error("get user failed", zap.Error(err), zap.Uint32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// -----------------------------------------------------------------------
// PUT /users/:id
// -----------------------------------------------------------------------

// UpdateUser godoc
// @Summary Update a user
// @Param id path int true "User ID"
// @Param body body models.UpdateUserRequest true "Update payload"
// @Success 200 {object} models.UserResponse
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if errs := h.validate.Struct(req); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{
			Error:   "validation failed",
			Details: formatValidationErrors(errs),
		})
	}

	resp, err := h.svc.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.log.Error("update user failed", zap.Error(err), zap.Uint32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// -----------------------------------------------------------------------
// DELETE /users/:id
// -----------------------------------------------------------------------

// DeleteUser godoc
// @Summary Delete a user
// @Param id path int true "User ID"
// @Success 204
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		h.log.Error("delete user failed", zap.Error(err), zap.Uint32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// -----------------------------------------------------------------------
// GET /users?page=1&limit=10
// -----------------------------------------------------------------------

// ListUsers godoc
// @Summary List users with pagination
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} models.ListUsersResponse
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page := int32(c.QueryInt("page", 1))
	limit := int32(c.QueryInt("limit", 10))

	resp, err := h.svc.List(c.Context(), page, limit)
	if err != nil {
		h.log.Error("list users failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------

func parseIDParam(c *fiber.Ctx) (uint32, error) {
	raw := c.Params("id")
	val, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

func formatValidationErrors(err error) map[string]string {
	out := make(map[string]string)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			out[fe.Field()] = fe.Tag()
		}
	}
	return out
}
