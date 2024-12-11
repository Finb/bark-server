package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"

	"github.com/mritd/logger"
)

type CommonResp struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type routerFunc struct {
	Name   string
	Weight int
	Func   func(router fiber.Router)
}

type routeSlice []routerFunc

func (r routeSlice) Len() int { return len(r) }

func (r routeSlice) Less(i, j int) bool { return r[i].Weight > r[j].Weight }

func (r routeSlice) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

var routerOnce sync.Once
var routes routeSlice

// register new route with key name
// key name is used to eliminate duplicate routes
// key name not case sensitive
func registerRoute(name string, f func(router fiber.Router)) {
	registerRouteWithWeight(name, 50, f)
}

// register new route with weight
func registerRouteWithWeight(name string, weight int, f func(router fiber.Router)) {
	if weight > 100 || weight < 0 {
		logger.Fatalf("route [%s] weight must be >= 0 and <=100", name)
	}

	for _, r := range routes {
		if strings.EqualFold(name, r.Name) {
			logger.Fatalf("route [%s] already registered", r.Name)
		}
	}

	routes = append(routes, routerFunc{
		Name:   name,
		Weight: weight,
		Func:   f,
	})
}

func routerSetup(router fiber.Router) {
	routerOnce.Do(func() {
		router.Use(fiberlogger.New(fiberlogger.Config{
			Format:     "${time}     INFO    ${ip} -> [${status}] ${method} ${latency} ${route} => ${url} ${body}\n",
			TimeFormat: "2006-01-02 15:04:05",
			Output:     os.Stdout,
		}))
		router.Use(fiberrecover.New())
		sort.Sort(routes)
		for _, r := range routes {
			r.Func(router)
			logger.Infof("load route [%s] success...", r.Name)
		}
	})
}

// for the fast return success result
func success() CommonResp {
	return CommonResp{
		Code:      200,
		Message:   "success",
		Timestamp: time.Now().Unix(),
	}
}

// for the fast return failed result
func failed(code int, message string, args ...interface{}) CommonResp {
	return CommonResp{
		Code:      code,
		Message:   fmt.Sprintf(message, args...),
		Timestamp: time.Now().Unix(),
	}
}

// for the fast return result with custom data
func data(data interface{}) CommonResp {
	return CommonResp{
		Code:      200,
		Message:   "success",
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}
