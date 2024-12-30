package main

import (
    "flag"
    "fmt"
    "github.com/joyant/transme/proxy"
    "os"
)

func main() {
    var (
        from        string
        to          string
        logFilePath string
    )
    flag.StringVar(&logFilePath, "config", "/tmp/transme.log", "log file path")
    flag.StringVar(&from, "from", "", "The source address for the proxy server to listen on (e.g., :8080).")
    flag.StringVar(&to, "to", "", "The target address for the proxy server to forward connections to (e.g., localhost:9090).")
    flag.Parse()
    if from == "" {
        fmt.Println("Error: The 'from' flag is required but not provided. Please specify the source address using the '-from' flag.")
        flag.Usage()
        os.Exit(1)
    }
    if to == "" {
        fmt.Println("Error: The 'to' flag is required but not provided. Please specify the target address using the '-to' flag.")
        flag.Usage()
        os.Exit(1)
    }
    err := os.Setenv("LOG_FILE_PATH", logFilePath)
    if err != nil {
        fmt.Printf("Error: Unable to set LOG_FILE_PATH environment variable. hint: %s", err.Error())
        flag.Usage()
        os.Exit(1)
    }
    fmt.Printf("from: %s to: %s logFilePath: %s\n", from, to, logFilePath)
    proxy.Start(from, to)
}
