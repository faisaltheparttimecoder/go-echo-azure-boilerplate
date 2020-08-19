package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func main() {
	// Kickoff by initializing the logger
	initLogger()

	// Start the router and serve the API traffic
	startWebServer()
}

// HTML Templates Render
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Start the web server to serve traffic
func startWebServer() {
	Debugf("Starting the echo instance and serving the API traffic")
	// Starting a new echo instance
	e := echo.New()

	// Middleware sessions
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionName))))

	// Middleware remove any "/" from URL, since most of the API we used doesn't have
	// trailing slack
	e.Pre(middleware.RemoveTrailingSlash())

	// Load all the public facing templates
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	// Error Handler to send all the internal error with a message
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, err)
		}
	}

	// Serve the routes
	webRouter(e)

	// Start server
	e.Logger.Fatal(e.Start(":" + IsSettingEmpty("API_PORT")))
}

// The router that is going to server the API traffic
func webRouter(e *echo.Echo) {
	// In & Out passage to the app
	e.GET("/login", loginHandler)
	e.GET("/auth/azure", azureLoginHandler)
	e.GET(IsSettingEmpty("AZURE_REDIRECT_URL"), azureCallbackHandler)
	e.GET("/auth/azure/token", azureTokenHandler)
	e.GET("/logout", logoutHandler)

	// Web Pages
	e.GET("/", homeHandler, AuthenticateRequestMiddleWare)
	e.GET("/restricted", restrictedHandler, AuthenticateRequestMiddleWare)
	e.GET("/unrestricted", unrestrictedHandler)
}
