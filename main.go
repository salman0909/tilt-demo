package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {

    // ===================================================================
    // App Starting
    log.Printf("main : Started")
    defer log.Println("main : Completed")
    // ===================================================================
    // Start HTTP server
    api := http.Server{
        Addr:         "localhost:8080",
        Handler:      http.HandlerFunc(ResponseHandler),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
    }
    // Make a channel to listen for errors coming from the listener.
    // Use a buffered channel so the goroutine can exit if we don't
    // collect this error.
    serverErrors := make(chan error, 1)
    // Start the service listening for requests.
    go func() {
        log.Printf("main : API listening on %s", api.Addr)
        serverErrors <- api.ListenAndServe()
    }()
    // Make a channel to listen for an interrupt or terminate signal
    // from the OS. Use a buffered channel because the signal package
    // requires it.
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
    // ===================================================================
    // Shutdown
    // Blocking main thread and waiting for shutdown.
    select {
    case err := <-serverErrors:
        log.Fatalf("error: listening and serving: %s", err)
    case <-shutdown:
        log.Println("main : Start shutdown")
        // Give outstanding requests a deadline for completion.
        const timeout = 5 * time.Second
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        // Asking listener to shutdown and load shed.
        err := api.Shutdown(ctx)
        if err != nil {
            log.Printf(
                "main : Graceful shutdown did not complete in %v : %v",
                timeout, err)
            err = api.Close()
        }
        if err != nil {
            log.Fatalf(
                "main : could not stop server gracefully : %v",
                err)
        }
    }
}

type Response struct {
    Message string `json:"message"`
}

// ResponseHandler is an HTTP Handler for returning a message.
func ResponseHandler(w http.ResponseWriter, r *http.Request) {

    responseJson := Response{"Hello, World!"}

    data, err := json.Marshal(responseJson)
    if err != nil {
        log.Println("error marshalling result", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    if _, err := w.Write(data); err != nil {
        log.Println("error writing result", err)
    }
}
