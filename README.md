# client-server-API
# Go Currency Exchange Rate Service  This project implements a client-server application in Go, showcasing skills in HTTP web servers, context management, database operations, and file handling.

## Project Overview

The project consists of two main components:

1. **Client (`client.go`)**: This component requests the current USD to BRL exchange rate from the server and saves the result to a text file.
2. **Server (`server.go`)**: This component handles incoming requests for exchange rates, fetches the data from an external API, stores the rates in a SQLite database, and returns the rate in JSON format to the client.

## Detailed Requirements

### Client (`client.go`)

- Makes an HTTP request to the server at the `/cotacao` endpoint to request the current USD to BRL exchange rate.
- Uses the `context` package to set a maximum timeout of 300ms for receiving the server's response.
- Logs an error if the timeout is exceeded.
- Saves the current exchange rate (field `bid` from the JSON response) to a file named `cotacao.txt` in the format: `Dólar: {value}`.

### Server (`server.go`)

- Listens for incoming HTTP requests at the `/cotacao` endpoint on port 8080.
- Uses the `context` package to set a maximum timeout of 200ms for calling the external exchange rate API.
- Fetches the exchange rate data from `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
- Returns the current exchange rate (field `bid` from the JSON response) to the client in JSON format.
- Logs an error if the API call exceeds the timeout.
- Uses the `context` package to set a maximum timeout of 10ms for saving the exchange rate data to a SQLite database.
- Logs an error if saving the data exceeds the timeout.

## Usage

### Setting Up

1. Ensure Go is installed on your machine.
2. Create a SQLite database named `cotacoes.db` in the same directory as the Go files.

### Running the Server

1. Navigate to the directory containing `server.go`.
2. Run the server using the command: `go run server.go`.

### Running the Client

1. Navigate to the directory containing `client.go`.
2. Run the client using the command: `go run client.go`.

The client will request the current exchange rate from the server and save it to `cotacao.txt`.
