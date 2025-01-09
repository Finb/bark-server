package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"io"

	"github.com/finb/bark-server/v2/database"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

// Before running the tests, a valid deviceToken must be set. Otherwise, the tests will fail.
const (
	deviceToken = ""
	key         = "MemoryBaseKey"
)

var app *fiber.App

func TestMain(m *testing.M) {
	if deviceToken == "" {
		panic("deviceToken is not set")
	}
	db = database.NewMemBase()
	db.SaveDeviceTokenByKey(key, deviceToken)
	app = NewServer()
	m.Run()
}

func TestRegister(t *testing.T) {
	Endpoint(t, []APITestCase{
		{
			Name:           "Normal registration",
			Method:         "GET",
			URL:            "/register?devicetoken=" + deviceToken,
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "Registration with key",
			Method:         "GET",
			URL:            "/register?key=" + key + "&devicetoken=" + deviceToken,
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "Registration with wrong key",
			Method:         "GET",
			URL:            "/register?key=" + "wrongKey" + "&devicetoken=" + deviceToken,
			Body:           "",
			IsJson:         false,
			WantStatusCode: 500,
		},
		{
			Name:           "Registration without devicetoken",
			Method:         "GET",
			URL:            "/register?key=" + key,
			Body:           "",
			IsJson:         false,
			WantStatusCode: 400,
		},
	})
}

func TestPushTitleAndBody(t *testing.T) {
	// Correct push
	Endpoint(t, []APITestCase{
		{
			Name:           "GET push body",
			Method:         "GET",
			URL:            "/" + key + "/body",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "GET push title body",
			Method:         "GET",
			URL:            "/" + key + "/title/body",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "GET push title subtitle body",
			Method:         "GET",
			URL:            "/" + key + "/title/subtitle/body",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST push body",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "body=body",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST push title body",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "title=title&body=body",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST push title subtitle body",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "title=title&subtitle=subtitle&body=body",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "GET title subtitle body URL parameters",
			Method:         "GET",
			URL:            "/" + key + "?title=title&subtitle=subtitle&body=body",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST title subtitle body POST parameters",
			Method:         "GET",
			URL:            "/" + key,
			Body:           "title=title&subtitle=subtitle&body=body",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST title subtitle body JSON parameters",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "{\"title\":\"title\",\"subtitle\":\"subtitle\",\"body\":\"body\"}",
			IsJson:         true,
			WantStatusCode: 200,
		},
		{
			Name:           "POST V2 title subtitle body",
			Method:         "POST",
			URL:            "/push",
			Body:           "device_key=" + key + "&title=title&subtitle=subtitle&body=body",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST title subtitle body JSON parameters V2",
			Method:         "POST",
			URL:            "/push",
			Body:           "{\"title\":\"title\",\"subtitle\":\"subtitle\",\"body\":\"body\",\"device_key\":\"" + key + "\"}",
			IsJson:         true,
			WantStatusCode: 200,
		},
	})

	// Incorrect push
	Endpoint(t, []APITestCase{
		{
			Name:           "GET push without key",
			Method:         "GET",
			URL:            "/body",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 400,
		},
		{
			Name:           "POST push without key",
			Method:         "POST",
			URL:            "/push",
			Body:           "title=title&subtitle=subtitle&body=body",
			IsJson:         false,
			WantStatusCode: 400,
		},
		{
			Name:           "POST JSON push without key",
			Method:         "POST",
			URL:            "/push",
			Body:           "body=body",
			IsJson:         true,
			WantStatusCode: 400,
		},
		{
			Name:           "GET push with too many parameters",
			Method:         "GET",
			URL:            "/" + key + "/title/subtitle/body/extra",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 404,
		},
	})
}

func TestCiphertext(t *testing.T) {
	Endpoint(t, []APITestCase{
		{
			Name:           "Send encrypted push",
			Method:         "GET",
			URL:            "/" + key + "/body?ciphertext=text&iv=01234567890123456",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "Send encrypted push, omit body",
			Method:         "GET",
			URL:            "/" + key + "?ciphertext=text",
			Body:           "",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST send encrypted push",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "ciphertext=text",
			IsJson:         false,
			WantStatusCode: 200,
		},
		{
			Name:           "POST send encrypted push V2",
			Method:         "POST",
			URL:            "/push",
			Body:           "{\"device_key\":\"" + key + "\",\"ciphertext\":\"text\"}",
			IsJson:         true,
			WantStatusCode: 200,
		},
	})
}

func TestBatchPush(t *testing.T) {
	Endpoint(t, []APITestCase{
		{
			Name:           "Batch Push",
			Method:         "POST",
			URL:            "/" + key,
			Body:           "{\"title\":\"title\",\"subtitle\":\"subtitle\",\"body\":\"body\",\"device_keys\":[\"" + key + "\",\"" + key + "\",\"" + key + "\"]}",
			IsJson:         true,
			WantStatusCode: 200,
		},
		{
			Name:           "Batch Push",
			Method:         "POST",
			URL:            "/push",
			Body:           "{\"title\":\"title\",\"subtitle\":\"subtitle\",\"body\":\"body\",\"device_keys\":[\"" + key + "\",\"" + key + "\",\"" + key + "\"]}",
			IsJson:         true,
			WantStatusCode: 200,
		},
		{
			Name:           "Batch Push",
			Method:         "POST",
			URL:            "/push",
			Body:           "{\"title\":\"title\",\"subtitle\":\"subtitle\",\"body\":\"body\",\"device_keys\": \"" + key + "," + key + "," + key + "\"}",
			IsJson:         true,
			WantStatusCode: 200,
		},
	})
}

type APITestCase struct {
	Name           string
	Method         string
	URL            string
	Body           string
	IsJson         bool
	WantStatusCode int
}

func NewServer() *fiber.App {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
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

	routerSetup(fiberApp)
	return fiberApp
}

func Endpoint(t *testing.T, tc []APITestCase) {
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.Method, tt.URL, bytes.NewBufferString(tt.Body))
			if tt.IsJson {
				req.Header.Set("Content-Type", "application/json")
			} else {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != tt.WantStatusCode {
				body, _ := io.ReadAll(io.Reader(res.Body))
				t.Fatalf("want %d, got %d, res: %s", tt.WantStatusCode, res.StatusCode, string(body))
			}
		})
		// Prevent rate limiting by sending requests too quickly
		time.Sleep(100 * time.Millisecond)
	}
}
