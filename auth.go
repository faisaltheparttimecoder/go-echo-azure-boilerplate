package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Custom Middleware to basically check if all the registered routes
// are authenticated i.e JWT Token verification & allow only if the token is valid
func AuthenticateRequestMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	Debugf("Verifying if the URL is authenticated")
	return func(c echo.Context) error {
		_, err := validateJwtToken(c)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		return next(c)
	}
}

// Validate the JWT token
func validateJwtToken(c echo.Context) (*jwt.Token, error){
	Debugf("Validating the JWT Token")
	var tk *jwt.Token

	// Get the JWT from session
	t, err := obtainJWTToken(c)
	if err != nil {
		return tk, err
	}

	// Read the JWT
	claims := jwt.MapClaims{}
	tk, err = jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return tk, PrintErrorAndReturn("Redirecting to login page due to JWT parsing error: %v", err)
	}

	// If the token valid
	if !tk.Valid {
		return tk, PrintErrorAndReturn("%v", errors.New("jwt token found is not valid"))
	}
	return tk, nil
}

// Get the session for jwt token
func obtainJWTToken(c echo.Context) (string, error) {
	Debug("Obtaining the session that keeps the JWT Token")
	sess, err := getSession(c)
	if err != nil {
		return "", PrintErrorAndReturn("Unable to get session from the browser: %v", err)
	}
	t := sess.Values[jwtName]
	if t == nil {
		return "", PrintErrorAndReturn("%v", errors.New("redirecting to login page since the session was empty"))
	}
	return t.(string), nil
}