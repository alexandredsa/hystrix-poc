package circuitbreaker

import (
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo"
)

var ConfOriginPort = ":7000"
var ConfBallbackPort = ":7001"
var OriginPort = 7000
var FallbackPort = 7001

func init() {
	go setupEndpoint(createHystrixHandler)
}

func setupEndpoint(handler func(c echo.Context) error) {
	e := echo.New()
	e.GET("/test", handler)
	e.Start(":7890")
}

func createThirdService(name string, httpCode int, interval time.Duration) *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		time.Sleep(interval)
		c.String(httpCode, name)
		return nil
	})
	return e
}

func createHystrixHandler(c echo.Context) error {
	hystrix.ConfigureCommand("test_endpoint", hystrix.CommandConfig{
		Timeout:               5000,
		MaxConcurrentRequests: 3,
		ErrorPercentThreshold: 10,
	})

	output := make(chan int, 1)
	errs := hystrix.Go("test_endpoint", func() error {
		resp, err := http.Get("http://127.0.0.1" + ConfOriginPort + "/")
		if err == nil {
			output <- resp.StatusCode
			return nil
		}
		return err
	}, func(err error) error {
		resp, err := http.Get("http://127.0.0.1" + ConfBallbackPort + "/")
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
