package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/finb/bark-server/v2/database"

	jsoniter "github.com/json-iterator/go"

	"github.com/gofiber/fiber/v2"

	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

var db database.Database

func main() {
	app := &cli.App{
		Name:    "bark-server",
		Usage:   "Push Server For Bark",
		Version: fmt.Sprintf("%s %s %s", version, commitID, buildDate),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server listen address",
				EnvVars: []string{"BARK_SERVER_ADDRESS"},
				Value:   "0.0.0.0:8080",
			},
			&cli.StringFlag{
				Name:    "url-prefix",
				Usage:   "Serve URL Prefix",
				EnvVars: []string{"BARK_SERVER_URL_PREFIX"},
				Value:   "/",
			},
			&cli.StringFlag{
				Name:    "data",
				Usage:   "Server data storage dir",
				EnvVars: []string{"BARK_SERVER_DATA_DIR"},
				Value:   "/data",
			},
			&cli.StringFlag{
				Name:    "dsn",
				Usage:   "MySQL DSN user:pass@tcp(host)/dbname",
				EnvVars: []string{"BARK_SERVER_DSN"},
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "serverless",
				Usage:   "serverless mode",
				EnvVars: []string{"BARK_SERVER_SERVERLESS"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "cert",
				Usage:   "Server TLS certificate",
				EnvVars: []string{"BARK_SERVER_CERT"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "key",
				Usage:   "Server TLS certificate key",
				EnvVars: []string{"BARK_SERVER_KEY"},
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "case-sensitive",
				Usage:   "Enable HTTP URL case sensitive",
				EnvVars: []string{"BARK_SERVER_CASE_SENSITIVE"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "strict-routing",
				Usage:   "Enable strict routing distinction",
				EnvVars: []string{"BARK_SERVER_STRICT_ROUTING"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "reduce-memory-usage",
				Usage:   "Aggressively reduces memory usage at the cost of higher CPU usage if set to true",
				EnvVars: []string{"BARK_SERVER_REDUCE_MEMORY_USAGE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "Basic auth username",
				EnvVars: []string{"BARK_SERVER_BASIC_AUTH_USER"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "Basic auth password",
				EnvVars: []string{"BARK_SERVER_BASIC_AUTH_PASSWORD"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "proxy-header",
				Usage:   "The remote IP address used by the bark server http header",
				EnvVars: []string{"BARK_SERVER_PROXY_HEADER"},
				Value:   "",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Usage:   "Maximum number of concurrent connections",
				EnvVars: []string{"BARK_SERVER_CONCURRENCY"},
				Value:   256 * 1024,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "read-timeout",
				Usage:   "The amount of time allowed to read the full request, including the body",
				EnvVars: []string{"BARK_SERVER_READ_TIMEOUT"},
				Value:   3 * time.Second,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "write-timeout",
				Usage:   "The maximum duration before timing out writes of the response",
				EnvVars: []string{"BARK_SERVER_WRITE_TIMEOUT"},
				Value:   3 * time.Second,
				Hidden:  true,
			},
			&cli.DurationFlag{
				Name:    "idle-timeout",
				Usage:   "The maximum amount of time to wait for the next request when keep-alive is enabled",
				EnvVars: []string{"BARK_SERVER_IDLE_TIMEOUT"},
				Value:   10 * time.Second,
				Hidden:  true,
			},
		},
		Authors: []*cli.Author{
			{Name: "mritd", Email: "mritd@linux.com"},
			{Name: "Finb", Email: "to@day.app"},
		},
		Action: func(c *cli.Context) error {
			fiberApp := fiber.New(fiber.Config{
				ServerHeader:      "Bark",
				CaseSensitive:     c.Bool("case-sensitive"),
				StrictRouting:     c.Bool("strict-routing"),
				Concurrency:       c.Int("concurrency"),
				ReadTimeout:       c.Duration("read-timeout"),
				WriteTimeout:      c.Duration("write-timeout"),
				IdleTimeout:       c.Duration("idle-timeout"),
				ProxyHeader:       c.String("proxy-header"),
				ReduceMemoryUsage: c.Bool("reduce-memory-usage"),
				JSONEncoder:       jsoniter.Marshal,
				Network:           "tcp",
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					code := fiber.StatusInternalServerError
					if e, ok := err.(*fiber.Error); ok {
						code = e.Code
					}
					return c.Status(code).JSON(CommonResp{
						Code:      code,
						Message:   err.Error(),
						Timestamp: time.Now().Unix(),
					})
				},
			})

			fiberRouter := fiberApp.Group(c.String("url-prefix"))
			routerAuth(c.String("user"), c.String("password"), fiberRouter)
			routerSetup(fiberRouter)

			if serverless := c.Bool("serverless"); serverless {
				// use system environment variable.
				db = database.NewEnvBase()
			} else if dsn := c.String("dsn"); dsn != "" {
				db = database.NewMySQL(dsn)
			} else {
				db = database.NewBboltdb(c.String("data"))
			}

			go func() {
				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
				for range sigs {
					logger.Warn("Received a termination signal, bark server shutdown...")
					if err := fiberApp.Shutdown(); err != nil {
						logger.Errorf("Server forced to shutdown error: %v", err)
					}
					if err := db.Close(); err != nil {
						logger.Errorf("Database close error: %v", err)
					}
				}
			}()

			logger.Infof("Bark Server Listen at: %s , Database: %s", c.String("addr"), reflect.TypeOf(db))
			if cert, key := c.String("cert"), c.String("key"); cert != "" && key != "" {
				return fiberApp.ListenTLS(c.String("addr"), cert, key)
			}
			return fiberApp.Listen(c.String("addr"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
