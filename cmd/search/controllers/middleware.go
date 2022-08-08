package controllers

import (
	"fmt"
	"net/http"

	"github.com/evan-buss/openbooks/cmd/search/model"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type UserInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

func setUserSession(request *http.Request, user *model.User) *sessions.Session {
	session, _ := store.Get(request, "user")
	// Set some session values.
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.Values["id"] = user.ID
	session.Values["name"] = user.Username

	return session
}

func getUserSession(request *http.Request) (*sessions.Session, error) {
	return store.Get(request, "user")
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := getUserSession(c.Request())

		if err != nil || session.IsNew {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized. Please sign in.")
		}

		id := session.Values["id"]
		name := session.Values["name"]
		fmt.Println(id, name)

		userInfo := &UserInfo{
			ID:   session.Values["id"].(uint),
			Name: session.Values["name"].(string),
		}

		c.Set("user", userInfo)
		return next(c)
	}
}
