package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wine-shop-api/internal/service"
)

type UploadHandler struct {
	CloudinaryService *service.CloudinaryService
}

// UploadImage godoc
// @Summary      Upload an image
// @Description  Upload an image to Cloudinary
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "Image file"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /admin/upload [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Upload to Cloudinary
	url, err := h.CloudinaryService.UploadImage(file, "wine-shop/products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"message": "Image uploaded successfully",
	})
}
