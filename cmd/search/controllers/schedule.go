package controllers

import (
	"net/http"
	"strconv"

	"github.com/evan-buss/openbooks/cmd/search/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) ScheduleRoutes() {
	s.Echo.GET("/api/schedule", s.GetSchedules())
	s.Echo.POST("/api/schedule", s.CreateSchedule())
	s.Echo.PUT("/api/schedule", s.UpdateSchedule())
	s.Echo.DELETE("/api/schedule/:id", s.DeleteSchedule())
}

// GetSchedules godoc
// @Summary Get all schedules
// @Tags schedule
// @Success 200 {array} model.Schedule
// @Failure 400 {object} controllers.HTTPError
// @Router /api/schedule [get]
func (s *Server) GetSchedules() echo.HandlerFunc {
	return func(e echo.Context) error {
		schedules, err := model.FindAllSchedules(s.DB, e.Request().Context())

		if err != nil {
			s.Logger.Error("error", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return e.JSON(http.StatusOK, schedules)
	}
}

// CreateSchedule godoc
// @Summary Create a schedule
// @Tags schedule
// @Param schedule body model.Schedule true "Schedule"
// @Success 201 {object} model.Schedule
// @Failure 400 {object} controllers.HTTPError
// @Router /api/schedule [post]
func (s *Server) CreateSchedule() echo.HandlerFunc {
	return func(e echo.Context) error {
		schedule := new(model.Schedule)
		if err := e.Bind(schedule); err != nil {
			s.Logger.Error("error", zap.Error(err))
			return err
		}

		err := model.CreateSchedule(s.DB, e.Request().Context(), schedule)
		if err != nil {
			s.Logger.Error("error", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		newSchedule, err := model.GetScheduleById(s.DB, e.Request().Context(), schedule.ID)
		if err != nil {
			s.Logger.Error("error", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return e.JSON(http.StatusCreated, newSchedule)
	}
}

// UpdateSchedule godoc
// @Summary Update a schedule
// @Tags schedule
// @Param schedule body model.Schedule true "Schedule"
// @Success 200 {object} model.Schedule
// @Failure 400 {object} controllers.HTTPError
// @Router /api/schedule [put]
func (s *Server) UpdateSchedule() echo.HandlerFunc {
	return func(e echo.Context) error {
		schedule := new(model.Schedule)
		if err := e.Bind(schedule); err != nil {
			s.Logger.Error("error", zap.Error(err))
			return err
		}
		err := model.UpdateSchedule(s.DB, e.Request().Context(), schedule)

		if err != nil {
			s.Logger.Error("error", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return e.JSON(http.StatusOK, schedule)
	}
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Tags schedule
// @Param scheduleId path int true "Schedule ID"
// @Success 204
// @Failure 400 {object} controllers.HTTPError
// @Router /api/schedule/{scheduleId} [delete]
func (s *Server) DeleteSchedule() echo.HandlerFunc {
	return func(e echo.Context) error {
		param := e.Param("id")

		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			s.Logger.Sugar().Infof("error converting id to int: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}

		err = model.DeleteSchedule(s.DB, e.Request().Context(), uint(id))

		if err != nil {
			s.Logger.Sugar().Errorf("error deleting schedule: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "error deleting schedule")
		}

		return e.NoContent(http.StatusNoContent)
	}
}
