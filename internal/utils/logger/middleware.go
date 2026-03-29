package logger

import (
	"strings"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	gLogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func Fiber() fiber.Handler {
	return fiberzap.New(fiberzap.Config{
		SkipURIs: []string{"/api/v1/swagger/*"},
		SkipResBody: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.OriginalURL(), "/api/v1/swagger")
		},
		Logger: zapLog,
		Fields: []string{
			"requestId", "method", "path",
			"status", "error",
			//"reqHeaders",
			"queryParams", "bytesReceived", "body",
			"bytesSent", "resBody",
			"ip", "ua", "latency",
		},
		FieldsFunc: nil,
		Messages:   []string{"Internal server error", "Invalid request", "Success"},
	})
}

//const FiberFormat = "${time} \u001B[35mHTTP\u001B[0m\t${method} ${path}\t[${status}] ${error}\t" + `{"ip": "${ip}", "latency": "${latency}", "method": "${method}", "url": "${path}", "status": ${status}, "error": "${error}"}` + "\n"

//const FiberFormat = "${time} ${method}\t${path}\t[${status}] ${error}\t" + `{"ip": "${ip}", "latency": "${latency}", "method": "${method}", "url": "${path}", "status": ${status}, "error": "${error}"}` + "\n"

// Gorm provides a middleware for Gorm logging
func Gorm() gLogger.Interface {
	l := zapgorm2.New(zapLog)
	l.LogMode(gLogger.Info)
	return l
}
