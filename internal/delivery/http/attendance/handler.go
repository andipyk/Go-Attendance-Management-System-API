package attendance

import (
	"net/http"
	"time"

	"golang-tes/internal/domain"
	"golang-tes/internal/utils"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	attendanceUsecase domain.AttendanceUsecase
}

func NewAttendanceHandler(attendanceUsecase domain.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceUsecase: attendanceUsecase,
	}
}

type markAttendanceRequest struct {
	Status string `json:"status" binding:"omitempty,oneof=present absent late"`
}

type getAttendanceRequest struct {
	Date string `form:"date" binding:"required" time_format:"2006-01-02"`
}

func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "user ID not found in context")
		return
	}

	var req markAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	attendance := &domain.Attendance{
		UserID: userID,
		Date:   time.Now(),
		Status: req.Status,
	}

	err := h.attendanceUsecase.MarkAttendance(c.Request.Context(), attendance)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark attendance", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Attendance marked successfully", nil)
}

func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	var req getAttendanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err.Error())
		return
	}

	attendances, err := h.attendanceUsecase.GetAttendanceByDate(c.Request.Context(), date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get attendance records", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Attendance records retrieved successfully", attendances)
}

func (h *AttendanceHandler) GetUserAttendance(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "user ID not found in context")
		return
	}

	attendances, err := h.attendanceUsecase.GetUserAttendance(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user attendance records", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User attendance records retrieved successfully", attendances)
}
