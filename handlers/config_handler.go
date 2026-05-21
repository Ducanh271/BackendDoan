package handlers

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	ConfigService *service.ConfigService
}

func NewConfigHandler(cs *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{ConfigService: cs}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	cfg, err := h.ConfigService.GetConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": cfg})
}

func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Dữ liệu không hợp lệ"})
		return
	}

	// Mặc định update cho cấu hình ID = 1
	if err := h.ConfigService.UpdateConfig(1, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cập nhật cấu hình thành công"})
}

func (h *ConfigHandler) AddWifi(c *gin.Context) {
	var req dto.AddWifiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Dữ liệu không hợp lệ"})
		return
	}

	if err := h.ConfigService.AddWifi(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Thêm Wi-Fi thành công"})
}

func (h *ConfigHandler) DeleteWifi(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.ConfigService.DeleteWifi(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Lỗi xóa Wi-Fi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Xóa Wi-Fi thành công"})
}
