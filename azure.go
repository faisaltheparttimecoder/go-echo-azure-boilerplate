package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"time"
)

var (
	sessionName = IsSettingEmpty("SESSION_NAME")
	jwtName =   IsSettingEmpty("JWT_NAME")
	jwtKey = IsSettingEmpty("JWT_TOKEN")
	expiryHours = ConvertStringToNumber(IsSettingEmpty("EXPIRY_DATE"))
	sessionExpiryDate = 60 * 60 * expiryHours
)

type jwtClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Users struct {
	Name       string
	GivenName  string
	FamilyName string
	Email      string
}

// The url to authenticate with Azure.
func azureAuthUrl() string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?response_type=token"+
		"&client_id=%s&redirect_uri=%s&domain_hint=%s&resource=%s",
		IsSettingEmpty("AZURE_TENANT_ID"), IsSettingEmpty("AZURE_CLIENT_ID"),
		fmt.Sprintf("%s%s", IsSettingEmpty("WEB_URL"), IsSettingEmpty("AZURE_REDIRECT_URL")),
		IsSettingEmpty("AZURE_DOMAIN"), IsSettingEmpty("AZURE_RESOURCE"))
}

// Extract the token and validate the user
func extractAzureToken(c echo.Context) error {
	Debugf("Extracting token from the URL")
	err := deleteSession(c)
	if err != nil {
		return fmt.Errorf("failed in deleting the saved session: %v", err)
	}

	accessToken := c.FormValue("access_token")
	if !DoesValueExists(accessToken) {
		return fmt.Errorf("unable to access any token from azure")
	}
	return parseToken(accessToken, c)
}

// parse the token & Extract the user Information
func parseToken(t string, c echo.Context) error {
	Debugf("Parsing the azure token")
	token, _, err := new(jwt.Parser).ParseUnverified(t, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("failed in parsing the token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return extractUser(claims, c)
	}
	return fmt.Errorf("failed in claim the jwt token: %v", err)
}

// Extract the user information & store in the database
func extractUser(j jwt.MapClaims, c echo.Context) error {
	Debugf("Extracting the user information from the token")
	u := &Users{
		Email:      j["unique_name"].(string),
		Name:       j["name"].(string),
		GivenName:  j["given_name"].(string),
		FamilyName: j["family_name"].(string),
	}

	// Record the entry of the user in the database if needed
	// the "u" has all authenticated user details

	// JWT Token
	err := createJWTToken(u, c)
	if err != nil {
		return fmt.Errorf("failed to create JWT Token: %v", err)
	}

	return nil
}

// Create our own JWT Token and Store it on session no need to save this
// on database
func createJWTToken(u *Users, c echo.Context) error {
	Debugf("Creating the JWT Token")
	// New JWT token
	claims := &jwtClaims{
		u.Name,
		u.Email,
		jwt.StandardClaims{
			ExpiresAt: jwtExpiry(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return fmt.Errorf("failed to sign jwt token: %v", err)
	}

	// Save session
	err = createSession(c, t)
	if err != nil {
		return fmt.Errorf("failed to create jwt token: %v", err)
	}

	return nil
}

// Create session information and store it has a cookie
func createSession(c echo.Context, t string) error {
	sess, _ := getSession(c)
	sess.Values[jwtName] = t
	sess.Options.MaxAge = sessionExpiryDate
	err := saveSession(sess, c)
	if err != nil {
		return fmt.Errorf("failed to create session, err: %v", err)
	}
	return err
}

// Generate expiry date for JWT Token
func jwtExpiry() int64 {
	return time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix()
}

// Get Session
func getSession(c echo.Context) (*sessions.Session, error) {
	Debugf("Getting the session information")
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return sess, PrintErrorAndReturn("failed in getting the session information: %v", err)
	}
	return sess, nil
}

// Save Session
func saveSession(sess *sessions.Session, c echo.Context) error {
	Debugf("Saving session cookie")
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("failed in saving the session: %v", err)
	}
	return nil
}

// Delete session
// https://github.com/gorilla/sessions/blob/03b6f63cc43ef9c7240a635a5e22b13180e822b8/store.go#L116
func deleteSession(c echo.Context) error {
	Debugf("Deleting session cookie")
	sess, _ := getSession(c)
	sess.Options.MaxAge = -1
	err := saveSession(sess, c)
	if err != nil {
		return fmt.Errorf("failed to delete session, err: %v", err)
	}
	return nil
}
