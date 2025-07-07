package films

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

// CreateFilm godoc
// @Summary Create a new film
// @Description Create a new film with the provided details
// @Tags films
// @Accept json
// @Produce json
// @Param film body Film true "Film data to create"
// @Success 201 {object} Film "Successfully created film"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films [post]
func (h *Handler) CreateFilm(c *gin.Context) {
	var newFilm Film
	if err := c.ShouldBindJSON(&newFilm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	newFilm.ID = id.String()
	newFilm.CreatedAt = time.Now()
	newFilm.UpdatedAt = newFilm.CreatedAt

	if err := h.storage.Create(c.Request.Context(), newFilm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create film card"})
		return
	}
	c.JSON(http.StatusCreated, newFilm)
}

// GetList godoc
// @Summary Get all films
// @Description Retrieve a list of all films with optional sorting
// @Tags films
// @Produce json
// @Param sort_by query string false "Field to sort by (title, rating, release_date)"
// @Success 200 {array} Film "List of films"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films [get]
func (h *Handler) GetList(c *gin.Context) {
	films, err := h.storage.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetListSort godoc
// @Summary Get sorted films list
// @Description Retrieve a list of films sorted by specified criteria
// @Tags films
// @Produce json
// @Param sort_by query string false "Field to sort by (title, rating, release_date)"
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Success 200 {array} Film "Sorted list of films"
// @Failure 400 {object} map[string]string "Invalid sort parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films/sorted [get]
func (h *Handler) GetListSort(c *gin.Context) {
	films, err := h.storage.FindAllSort(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetUserFilm godoc
// @Summary Get films by user ID
// @Description Retrieve all films associated with specific user
// @Tags films
// @Produce json
// @Param uuid path string true "User ID (UUID)"
// @Success 200 {array} Film "List of user's films"
// @Success 200 {object} map[string]string "No films found for user"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films/user/{uuid} [get]
func (h *Handler) GetUserFilm(c *gin.Context) {
	param := c.Param("uuid")
	user, err := h.storage.FindOne(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find films"})
		return
	}

	if len(user) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No films found for this user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// PartiallyUpdateFilm godoc
// @Summary Partially update film
// @Description Update specific fields of a film
// @Tags films
// @Accept json
// @Produce json
// @Param uuid path string true "Film ID (UUID)"
// @Param updates body UpdateFilm true "Fields to update"
// @Success 204 "Film updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Film not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films/{uuid} [patch]
func (h *Handler) PartiallyUpdateFilm(c *gin.Context) {
	param := c.Param("uuid")
	var input UpdateFilm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.storage.PartialUpdate(c.Request.Context(), param, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update film card"})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteFilm godoc
// @Summary Delete a film
// @Description Remove a film by its UUID
// @Tags films
// @Produce json
// @Param uuid path string true "Film ID (UUID)"
// @Success 204 "Film deleted successfully"
// @Failure 404 {object} map[string]string "Film not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /films/{uuid} [delete]
func (h *Handler) DeleteFilm(c *gin.Context) {
	par := c.Param("uuid")
	if err := h.storage.Delete(c.Request.Context(), par); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete film"})
		return
	}
	c.Status(http.StatusNoContent)
}
