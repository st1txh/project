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

func (h *Handler) GetList(c *gin.Context) {
	users, err := h.storage.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

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

func (h *Handler) PartiallyUpdateUser(c *gin.Context) {
	param := c.Param("uuid")
	var input UserUpdate
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

func (h *Handler) DeleteUser(c *gin.Context) {
	par := c.Param("uuid")
	if err := h.storage.Delete(c.Request.Context(), par); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}
