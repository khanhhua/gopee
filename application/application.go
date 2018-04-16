package application

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SermoDigital/jose/jws"
	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"github.com/khanhhua/gopee/dao"
	"github.com/khanhhua/gopee/handlers"
	"github.com/khanhhua/gopee/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))
	if dao, dberr := dao.New(config.Get("CLEARDB_DATABASE_URL").(string)); dberr != nil {
		return nil, fmt.Errorf("Database connection failed")
	} else {
		app.dao = dao
	}

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dao          *dao.DAO
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetSessionStore(app.sessionStore))
	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "dao", app.dao))

			next.ServeHTTP(res, req)
		})
	})
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.Handle("/home", http.HandlerFunc(handlers.GetHome)).Methods("GET")
	// router.Handle("/console", http.HandlerFunc(handlers.ViewConsole))
	router.Handle("/call/{fnName}", http.HandlerFunc(handlers.Call)).Methods("POST")

	router.Handle("/auth/token", http.HandlerFunc(handlers.GetToken)).Methods("POST")
	router.Handle("/auth/dropbox", http.HandlerFunc(handlers.Authorize)).Methods("GET")

	// API routes are protected
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if jwt, err := jws.ParseJWTFromRequest(r); err != nil {
				http.Error(w, "Authentication required", 403)
			} else {
				claims := jwt.Claims()
				if sub, ok := claims.Subject(); ok != true {
					http.Error(w, "Authentication required", 403)
					return
				} else {
					logrus.Infoln("CLientKey: " + strings.Repeat("x", len(sub)))
					ctx := context.WithValue(r.Context(), "x-client-key", sub)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		})
	})
	apiRouter.Handle("/funs", http.HandlerFunc(handlers.Get)).Methods("GET")
	apiRouter.Handle("/funs/{id}", http.HandlerFunc(handlers.GetOne)).Methods("GET")
	apiRouter.Handle("/funs/{id}", http.HandlerFunc(handlers.Edit)).Methods("PUT")
	apiRouter.Handle("/funs", http.HandlerFunc(handlers.Compose)).Methods("POST")

	// Path of static files must be last!
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("static/console")))
	router.PathPrefix("/").HandlerFunc(indexHandler("static/console/index.html"))

	return router
}

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
