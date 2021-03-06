/**
 * File              : routes.go
 * Author            : Jiang Yitao <jiangyt.cn#gmail.com>
 * Date              : 11.08.2019
 * Last Modified Date: 08.09.2019
 * Last Modified By  : Jiang Yitao <jiangyt.cn#gmail.com>
 */
package api

import (
	"context"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	. "github.com/jiangytcn/gosms-ng/dash"
	. "github.com/jiangytcn/gosms-ng/logger"
	"github.com/jiangytcn/gosms-ng/pkg/api/config"
	"github.com/jiangytcn/gosms-ng/pkg/api/device"
	"github.com/jiangytcn/gosms-ng/pkg/api/sms"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func Routes() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	//r.Use(utils.Logger(logger))
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(middleware.Timeout(60 * time.Second))

	htmlData, err := Asset("index.html")
	if err != nil {
		panic(err)
		// Asset was not found.
	}

	tpl, _ := template.New("index").Parse(string(htmlData))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tpl.Execute(w, nil)
	})

	r.Get("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fs := http.StripPrefix("/static", http.FileServer(AssetFile()))
		fs := http.FileServer(AssetFile())

		fs.ServeHTTP(w, r)
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Use(apiVersionCtx("v1"))
		r.Mount("/smss", sms.Routes())
		r.Mount("/admin", config.Routes())
		r.Mount("/devices", device.Routes())
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		Logger.Info("route", zap.String("method", method), zap.String("route", route)) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		Logger.Fatal("Logging err", zap.Error(err)) // panic if there is an error
	}
	return r
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}
