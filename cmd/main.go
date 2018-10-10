package main

import (
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo"
)

func main() {
	// setupEndpoint(createHystrixHandler)
	setupThirdService(":1111", 202)
}

func setupEndpoint(handler func(c echo.Context) error) {
	e := echo.New()
	e.GET("/test", handler)
	e.Start(":7890")
}

func setupThirdService(port string, httpCode int) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		e.Logger.Info("Get on port " + port)
		c.String(httpCode, "")
		return nil
	})
	e.Logger.Info(e.Start(port))
}

func createHystrixHandler(c echo.Context) error {
	hystrix.ConfigureCommand("test_endpoint", hystrix.CommandConfig{
		Timeout:               5000,
		MaxConcurrentRequests: 3,
		ErrorPercentThreshold: 10,
	})

	output := make(chan int, 1)
	errs := hystrix.Go("test_endpoint", func() error {
		resp, err := http.Get("http://127.0.0.1:1111/")
		if err == nil {
			output <- resp.StatusCode
			return nil
		}
		return err
	}, func(err error) error {
		resp, err := http.Get("http://127.0.0.1:1313/")
		if err == nil {
			output <- resp.StatusCode
			return nil
		}
		return err
	})

	select {
	case out := <-output:
		c.String(out, "")
		return nil
	case err := <-errs:
		c.String(http.StatusGatewayTimeout, err.Error())
		return err
	}
}
