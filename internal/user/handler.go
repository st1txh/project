package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"rest-api-tutorial/pkg/logging"
	"time"
)

type Handler struct {
	logger  *logging.Logger
	storage *Storage
}

func NewHandler(storage *Storage, logger *logging.Logger) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data to create"
// @Success 201 {object} User "Successfully created user"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	newUser.ID = id.String()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = newUser.CreatedAt

	if err := h.storage.Create(c.Request.Context(), newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

// GetList godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func (h *Handler) GetList(c *gin.Context) {
	users, err := h.storage.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a single user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User ID (UUID)"
// @Success 200 {object} User "Requested user"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uuid} [get]
func (h *Handler) GetUser(c *gin.Context) {
	param := c.Param("uuid")
	user, err := h.storage.FindOne(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Fully update a user
// @Description Replace all user data with the provided values
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User ID (UUID)"
// @Param user body User true "Updated user data"
// @Success 204 "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uuid} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	param := c.Param("uuid")
	var (
		input User
	)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	input.ID = param
	input.UpdatedAt = time.Now()

	if err := h.storage.Update(c.Request.Context(), param, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.Status(http.StatusNoContent)
}

// PartiallyUpdateUser godoc
// @Summary Partially update a user
// @Description Update specific fields of a user
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User ID (UUID)"
// @Param updates body Update true "Fields to update"
// @Success 204 "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uuid} [patch]
func (h *Handler) PartiallyUpdateUser(c *gin.Context) {
	param := c.Param("uuid")
	var input Update
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.storage.PartialUpdate(c.Request.Context(), param, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User ID (UUID)"
// @Success 204 "User deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uuid} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	par := c.Param("uuid")
	if err := h.storage.Delete(c.Request.Context(), par); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}
