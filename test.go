package main

import (
    "context"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
)

// Mock server for testing
func TestServer(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
        defer cancel()

        select {
        case <-ctx.Done():
            http.Error(w, "Request timed out", http.StatusRequestTimeout)
            return
        case <-time.After(100 * time.Millisecond): // Simulate API delay
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(`{"USDBRL":{"bid":"5.42"}}`))
        }
    }))
    defer server.Close()

    client := &http.Client{}
    req, err := http.NewRequest("GET", server.URL+"/cotacao", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status OK but got %v", resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

    expected := `{"USDBRL":{"bid":"5.42"}}`
    if string(body) != expected {
        t.Fatalf("Expected body %s but got %s", expected, string(body))
    }
}

// Test client writing quote to file
func TestClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"bid":"5.42"}`))
    }))
    defer server.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/cotacao", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status OK but got %v", resp.Status)
    }

    var quote Quote
    if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    err = ioutil.WriteFile("quote.txt", []byte("Dólar: "+quote.Bid), 0644)
    if err != nil {
        t.Fatalf("Failed to write quote to file: %v", err)
    }

    data, err := ioutil.ReadFile("quote.txt")
    if err != nil {
        t.Fatalf("Failed to read quote file: %v", err)
    }

    expected := "Dólar: 5.42"
    if string(data) != expected {
        t.Fatalf("Expected %s but got %s", expected, string(data))
    }

    // Clean up
    os.Remove("quote.txt")
}
