package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Handler: Login | Handles the login part of the app
func loginHandler(c echo.Context) error {
	Debugf("Publishing the login page")
	return c.Render(http.StatusUnauthorized, "login.html", "")
}

// Handler: Azure Login | Send the page to azure authentication
func azureLoginHandler(c echo.Context) error {
	Debugf("Sending the page to azure for authentication")
	return c.Redirect(http.StatusTemporaryRedirect, azureAuthUrl())
}

// Handler: Login CallBack | Handles the callback after a successful authentication
func azureCallbackHandler(c echo.Context) error {
	Debugf("Request came back from azure, handling the request")
	return c.Render(http.StatusOK, "callback.html", "")
}

// Handler: Azure Token | Once successfully logged in its time to handle the toke generated by azure
func azureTokenHandler(c echo.Context) error {
	Debugf("Handling the token send by Azure")
	err := extractAzureToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("%v", err))
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Handler: Logout | Handles the logout part of the app
func logoutHandler(c echo.Context) error {
	Debugf("Trying to the log the user out")
	err := deleteSession(c)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to cleanup the session during logout: %v", err))
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// Handler: Home | Handler send to home page
func homeHandler(c echo.Context) error {
	Debugf("Publishing the home page")
	return c.Render(http.StatusUnauthorized, "home.html", "")
}

// Handler: Restricted | Handler send to restricted page
func restrictedHandler(c echo.Context) error {
	Debugf("Publishing the home page")
	t, _ := validateJwtToken(c)
	return c.Render(http.StatusUnauthorized, "restricted.html", t.Claims)
}

// Handler: UnRestricted | Handles send to unrestricted page
func unrestrictedHandler(c echo.Context) error {
	Debugf("Publishing the home page")
	return c.Render(http.StatusUnauthorized, "unrestricted.html", "")
}
