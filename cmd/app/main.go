package main

import (
	_ "embed"
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/makesalekz/gateway/internal/conf"
	"github.com/makesalekz/gateway/internal/server"

	_ "go.uber.org/automaxprocs"
)

//go:embed openapi.yaml
var openapiSpec []byte

var (
	Name    string = "gateway"
	Version string = "0.1.0"

	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.local.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *server.HTTPServer) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(hs),
	)
}

func main() {
	flag.Parse()
	logger := log.With(
		log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	bc, err := conf.Load(flagconf)
	if err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc, logger, openapiSpec)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
