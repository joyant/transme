package proxy

import (
    "github.com/joyant/transme/log"
    "go.uber.org/zap"
    "io"
    "net"
)

func proxy(source net.Conn, target net.Conn) {
    defer source.Close()
    defer target.Close()

    go func() {
        _, err := io.Copy(target, source)
        if err != nil {
            log.Logger.Error("error while proxying source to target", zap.Error(err))
        }
    }()
    _, err := io.Copy(source, target)
    if err != nil {
        log.Logger.Error("error while proxying target to source", zap.Error(err))
    }
}

func Start(from, to string) {
    listener, err := net.Listen("tcp", from)
    if err != nil {
        log.Logger.Error("failed to listen on", zap.String("from", from), zap.Error(err))
        return
    }
    defer listener.Close()

    log.Logger.Info("proxy server listening", zap.String("from", from), zap.String("to", to))

    for {
        fromConn, err := listener.Accept()
        if err != nil {
            log.Logger.Error("failed to accept client connection", zap.Error(err))
            continue
        }

        toConn, err := net.Dial("tcp", to)
        if err != nil {
            log.Logger.Error("failed to connect", zap.String("to", to), zap.Error(err))
            fromConn.Close()
            continue
        }

        log.Logger.Info("established connection")

        go proxy(fromConn, toConn)
    }
}
