package main

import (
    "context"
    "crypto/tls"
    "ghost2gemini/handlers"
    "ghost2gemini/models/config"
    "github.com/ninedraft/gemax/gemax"
    "log"
)

func main() {
    config.LoadConfig()

    server := &gemax.Server{
        Addr:    config.GeminiHost + ":" + config.GeminiPort,
        Handler: handlers.BlogHandler,
    }

    cert, errCert := tls.LoadX509KeyPair("data/cert/cert.pem", "data/cert/key.pem")
    if errCert != nil {
        log.Fatal(errCert)
    }

    cfg := &tls.Config{
        MinVersion:   tls.VersionTLS12,
        Certificates: []tls.Certificate{cert},
    }

    log.Println("Starting server listening on " + config.GeminiHost + ":" + config.GeminiPort)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    err := server.ListenAndServe(ctx, cfg)
    if err != nil {
        log.Printf("test server: Serve: %v", err)
    }
}
