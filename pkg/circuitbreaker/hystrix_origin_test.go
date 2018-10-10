package circuitbreaker

import (
	"net/http"
	"testing"
	"time"
)

func TestShouldReturnWithOrigin(t *testing.T) {
	e := createThirdService("ORIGIN", http.StatusOK, 1*time.Millisecond)
	go e.Start(ConfOriginPort)

	resp, err := http.Get("http://127.0.0.1:7890/test")

	if err != nil {
		t.Fatalf("Expected no errors")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %v but got %v", http.StatusOK, resp.StatusCode)
	}

	e.Close()
}
