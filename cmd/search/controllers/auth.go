package controllers

import (
	"net/http"

	"github.com/evan-buss/openbooks/cmd/search/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (api *Server) AuthRoutes() {
	api.POST("/api/auth/login", api.LogIn())
	api.POST("/api/auth/logout", api.LogOut())

	authGroup := api.Group("")
	authGroup.Use(IsAuthenticated)
	authGroup.GET("api/auth/me", api.Me())
}

type LogInRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin"`
}

// LogIn godoc
// @Summary      	Log in to the admin site
// @Description  	Requires a valid username and password. If first time setup, use admin / admin.
// @Id 				LogIn
// @Tags         	auth
// @Param  			credentials body LogInRequest true "Credentials"
// @Success      	200
// @Failure      	400
// @Router       /api/auth/login [post]
func (api *Server) LogIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &LogInRequest{}
		if err := c.Bind(request); err != nil {
			return err
		}

		user, err := model.FindUserByUsername(api.DB, c.Request().Context(), request.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid username or password.")
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid username or password.")
		}

		session := setUserSession(c.Request(), &user)

		// Save it before we write to the response/return from the handler.
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error saving session")
		}

		return c.NoContent(http.StatusOK)
	}
}

// LogOut godoc
// @Summary      	Log out of the admin site
// @Id 				LogOut
// @Tags         	auth
// @Success      	200
// @Failure      	401 {object} HTTPError "Unauthorized"
// @Router       	/api/auth/logout [post]
func (api *Server) LogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := getUserSession(c.Request())
		if err == nil {
			session.Options.MaxAge = -1
			session.Save(c.Request(), c.Response().Writer)
		}
		return c.NoContent(http.StatusOK)
	}
}

// Me godoc
// @Summary      	Get logged in user's information.
// @Id 				Me
// @Tags         	auth
// @Success      	200 {object} UserInfo "User information"
// @Failure      	401 {object} HTTPError "Unauthorized"
// @Router       	/api/auth/me [get]
func (api *Server) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		info := c.Get("user").(*UserInfo)
		return c.JSON(http.StatusOK, info)
	}
}
