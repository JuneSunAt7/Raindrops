package client

import (
    "context"
    "fmt"
    "net/http"
    "sync"
	"time"
)

var UserResponse string

var server *http.Server          
var serverRunning bool          
var serverMutex sync.Mutex        

func handleAction(w http.ResponseWriter, r *http.Request) {
    action := r.URL.Path[1:] // del first slash
    if action == "yes" {
        UserResponse = "yes"
    } else if action == "no" {
        UserResponse = "no"
    }
    w.WriteHeader(http.StatusOK)
}

func StartServer() {
    serverMutex.Lock()
    defer serverMutex.Unlock()

    if serverRunning {
        return 
    }

    server = &http.Server{Addr: ":8081"}

    http.HandleFunc("/", handleAction)

    go func() {
        err := server.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            fmt.Printf("Ошибка при запуске сервера: %v\n", err)
        }
    }()

    serverRunning = true
}

func StopServer() {
    serverMutex.Lock()
    defer serverMutex.Unlock()

    if !serverRunning {
        return 
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Ошибка при остановке сервера: %v\n", err)
    }

    serverRunning = false
}