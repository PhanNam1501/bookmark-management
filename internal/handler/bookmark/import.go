package bookmark

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler/dto"
	"github.com/PhanNam1501/bookmark-management/internal/handler/utils"
	"github.com/PhanNam1501/bookmark-management/internal/handler/validator"
	"github.com/PhanNam1501/bookmark-management/internal/service/queue"
	"github.com/PhanNam1501/bookmark-management/pkg/csv"
	"github.com/gin-gonic/gin"
)

const maxFileSize = 10 << 20 // 10MB

// Import godoc
// @Summary Bulk import bookmarks from CSV file
// @Description Upload a CSV file to asynchronously import multiple bookmarks for the authenticated user. The import job is queued and processed in the background.
// @Tags Bookmarks
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "CSV file with bookmark data"
// @Success 200 {object} dto.SuccessResponseString "Import job queued successfully"
// @Failure 400 {object} map[string]string "Invalid file or CSV parsing error"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bookmarks/import [post]
func (h *handler) Import(c *gin.Context) {
	// get uuid
	uid, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	// get .csv file from input
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	if file.Size > maxFileSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit (10MB)"})
		return
	}

	var data []*queue.ImportBookmarkInput
	err = csv.ParseFromMultipartFile(file, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid parse file to data"})
		return
	}
	// validate input from .csv
	if err := validator.ValidateImportData(data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create message and send message to queue
	err = h.q.SendImportBookmarkJob(c, uid, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// return OK
	c.JSON(http.StatusOK, &dto.SuccessResponse[string]{
		Message: "Import bookmarks successfully!",
	})
}
