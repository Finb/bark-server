package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/mritd/logger"
)

type CommonResp struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type routerFunc struct {
	Name   string
	Weight int
	Func   func(router *gin.Engine)
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
func registerRoute(name string, f func(router *gin.Engine)) {
	registerRouteWithWeight(name, 50, f)
}

// register new route with weight
func registerRouteWithWeight(name string, weight int, f func(router *gin.Engine)) {
	if weight > 100 || weight < 0 {
		logger.Fatalf("route [%s] weight must be >= 0 and <=100", name)
	}

	for _, r := range routes {
		if strings.ToLower(name) == strings.ToLower(r.Name) {
			logger.Fatalf("route [%s] already registered", r.Name)
		}
	}

	routes = append(routes, routerFunc{
		Name:   name,
		Weight: weight,
		Func:   f,
	})
}

var router *gin.Engine

func routerSetup(debug bool) {
	routerOnce.Do(func() {
		if debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router = gin.New()

		sort.Sort(routes)
		for _, r := range routes {
			r.Func(router)
			logrus.Infof("load route [%s] success...", r.Name)
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
