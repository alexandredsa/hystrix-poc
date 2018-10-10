package circuitbreaker

import (
	"net/http"
	"testing"
	"time"
)

func TestShouldReturnWithFallback(t *testing.T) {
	e := createThirdService("FALLBACK", http.StatusNoContent, 1*time.Millisecond)
	go e.Start(ConfBallbackPort)

	resp, err := http.Get("http://127.0.0.1:7890/test")

	if err != nil {
		t.Fatalf("Expected no errors")
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected %v but got %v", http.StatusNoContent, resp.StatusCode)
	}

	e.Close()
}
