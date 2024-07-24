package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type Quote struct {
    Bid string `json:"bid"`
}

func main() {
    db, err := sql.Open("sqlite3", "./quotes.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        bid TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error fetching quote from API: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        var result map[string]Quote
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            log.Printf("Error decoding API response: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        quote := result["USDBRL"]
        json.NewEncoder(w).Encode(quote)

        saveCtx, saveCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
        defer saveCancel()

        _, err = db.ExecContext(saveCtx, "INSERT INTO quotes (bid) VALUES (?)", quote.Bid)
        if err != nil {
            log.Printf("Error saving quote to database: %v", err)
        }
    })

    log.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
