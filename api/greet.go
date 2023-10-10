package main

import (
	_ "embed"

	"flag"
	"fmt"
	"net/http"
	"strings"

	"greet/api/internal/config"
	"greet/api/internal/handler"
	"greet/api/internal/svc"

	"github.com/swaggest/swgui/v5emb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

//go:embed greet.json
var spec []byte

var configFile = flag.String("f", "etc/greet.yaml", "the config file")

const (
	swaggerPath = "/api/doc"
	JsonPath    = "/api/greet.json"
	Titile      = "title"
)

var swaggerHandle http.Handler = v5emb.New(
	Titile,
	JsonPath,
	swaggerPath,
)

func Notfound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, swaggerPath) {
			swaggerHandle.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithNotFoundHandler(Notfound()))
	// or just like below
	// server = rest.MustNewServer(c.RestConf, rest.WithNotFoundHandler(swaggerHandle))

	defer server.Stop()
	// swagger  json file
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   JsonPath,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(spec)
		},
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Println("doc: http://localhost:8888/api/doc")
	server.Start()
}
