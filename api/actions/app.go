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
		auth.GET("/orders", OrdersList)                   // done
		auth.GET("/orders/{id}", OrdersIndex)             // done
		auth.PUT("/orders/{id}", OrdersUpdate)            // done
		auth.DELETE("/categories/{id}", CategoriesDelete) // done
		auth.POST("/categories", CategoriesCreate)        // done
		auth.PUT("/categories/{id}", CategoriesUpdate)    // done
		auth.DELETE("/items/{id}", ItemsDelete)           // done
		auth.POST("/items", ItemsCreate)                  // done
		auth.PUT("/items/{id}", ItemsUpdate)              // done
		// POST /items/{itemID}/picture - add picture to the item
		// POST /import/categories - upload list of categories (json)
		// POST /import/items - upload list of items (json)
		// POST /import/pictures - upload archive of pictures for existing items (zip, alias suffix matching)
		// GET /export/orders - download list of orders

		app.GET("/", HomeHandler)

		app.GET("/items", ItemsList)                 // done
		app.GET("/items/{id}", ItemsIndex)           // done
		app.GET("/categories", CategoriesList)       // done
		app.GET("/categories/{id}", CategoriesIndex) // done
		app.POST("/orders", OrdersCreate)            // done
		// PUT /orders/{orderID}/item - add item to the order (by id and count)

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
