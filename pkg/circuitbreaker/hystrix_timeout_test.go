package circuitbreaker

import (
	"net/http"
	"testing"
	"time"
)

/**
*
* The Hystrix is configured to wait 5 seconds per request
*
**/
func TestShouldReturnTimeoutWithoutFallback(t *testing.T) {
	e := createThirdService("ORIGIN", http.StatusOK, 10*time.Second)
	go e.Start(ConfOriginPort)

	resp, _ := http.Get("http://127.0.0.1:7890/test")

	if resp.StatusCode != http.StatusGatewayTimeout {
		t.Fatalf("Expected %v but got %v", http.StatusGatewayTimeout, resp.StatusCode)
	}

	e.Close()
}

func TestShouldReturnFallbackAfterTimeout(t *testing.T) {
	originEcho := createThirdService("ORIGIN", http.StatusOK, 30*time.Second)
	fallbackEcho := createThirdService("FALLBACK", http.StatusNoContent, 1*time.Second)
	go originEcho.Start(ConfOriginPort)
	go fallbackEcho.Start(ConfBallbackPort)

	resp, _ := http.Get("http://127.0.0.1:7890/test")

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected %v but got %v", http.StatusNoContent, resp.StatusCode)
	}

	originEcho.Close()
	fallbackEcho.Close()
}
