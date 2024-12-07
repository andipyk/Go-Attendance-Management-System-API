package attendance

import (
	"net/http"
	"time"

	"golang-tes/internal/domain"
	"golang-tes/internal/utils"
	"golang-tes/internal/utils/validator"

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

// MarkAttendance godoc
// @Summary Mark attendance
// @Description Mark attendance for the authenticated user
// @Tags attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body markAttendanceRequest true "Attendance status"
// @Success 201 {object} utils.Response "Attendance marked successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 409 {object} utils.Response "Attendance already marked"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /attendance [post]
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized.Error())
		return
	}

	var req markAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", domain.ErrInvalidInput.Error())
		return
	}

	if err := validator.ValidateAttendanceStatus(req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid attendance status", err.Error())
		return
	}

	attendance := &domain.Attendance{
		UserID: userID,
		Date:   time.Now(),
		Status: req.Status,
	}

	err := h.attendanceUsecase.MarkAttendance(c.Request.Context(), attendance)
	if err == domain.ErrAttendanceAlreadyMarked {
		utils.ErrorResponse(c, http.StatusConflict, "Failed to mark attendance", err.Error())
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark attendance", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Attendance marked successfully", nil)
}

// GetAttendance godoc
// @Summary Get all attendance records by date
// @Description Get attendance records for all users on a specific date
// @Tags attendance
// @Produce json
// @Security BearerAuth
// @Param date query string true "Date in YYYY-MM-DD format" Format(date)
// @Success 200 {object} utils.Response{data=[]domain.Attendance} "Attendance records retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /attendance [get]
func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	var req getAttendanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", domain.ErrInvalidInput.Error())
		return
	}

	date, err := time.Parse(domain.DateFormat, req.Date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", domain.ErrInvalidInput.Error())
		return
	}

	attendances, err := h.attendanceUsecase.GetAttendanceByDate(c.Request.Context(), date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get attendance records", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Attendance records retrieved successfully", attendances)
}

// GetUserAttendance godoc
// @Summary Get user attendance records
// @Description Get all attendance records for the authenticated user
// @Tags attendance
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.Attendance} "User attendance records retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /attendance/user [get]
func (h *AttendanceHandler) GetUserAttendance(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized.Error())
		return
	}

	attendances, err := h.attendanceUsecase.GetUserAttendance(c.Request.Context(), userID)
	if err == domain.ErrUserNotFound {
		utils.ErrorResponse(c, http.StatusNotFound, "Failed to get user attendance records", err.Error())
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user attendance records", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User attendance records retrieved successfully", attendances)
}
