package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/basicauth"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"goshop/api/models"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_api_session",
		})
		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		app.Use(middleware.PopTransaction(models.DB))

		// admin zone
		auth := app.Group("/")
		auth.Use(basicauth.Middleware(BasicAuth))
		auth.GET("/orders", OrdersList)
		auth.GET("/orders/{id}", OrdersIndex)
		auth.PUT("/orders/{id}", OrdersUpdate)
		auth.DELETE("/categories/{id}", CategoriesDelete)
		auth.POST("/categories", CategoriesCreate)
		auth.PUT("/categories/{id}", CategoriesUpdate)
		auth.DELETE("/items/{id}", ItemsDelete)
		auth.POST("/items", ItemsCreate)
		auth.PUT("/items/{id}", ItemsUpdate)
		auth.GET("/export/orders", ExportOrders)
		// auth.POST("/items/{id}/picture", ItemsAddPicture)
		auth.POST("/import/categories", ImportCategories) //doesn't work properly
		auth.POST("/import/items", ImportItems)           //doesn't work properly
		// POST /import/pictures - upload archive of pictures for existing items (zip, alias suffix matching)

		app.GET("/", HomeHandler)

		app.GET("/items", ItemsList)
		app.GET("/items/{id}", ItemsIndex)
		app.GET("/categories", CategoriesList)
		app.GET("/categories/{id}", CategoriesIndex)
		app.POST("/orders", OrdersCreate)
		app.PUT("/orders/{id}/item", AddItemToTheOrder)

		app.ServeFiles("/assets", assetBox)
	}

	return app
}

func BasicAuth(c buffalo.Context, basUsr string, basPwd string) (bool, error) {
	username := envy.Get("BASIC_USERNAME", "development")
	pwd := envy.Get("BASIC_PASSWORD", "development")
	if username == basUsr && pwd == basPwd {
		return true, nil
	}
	return false, basicauth.ErrAuthFail
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return ssl.ForceSSL(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
