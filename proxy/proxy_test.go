package proxy

import "testing"

func Test_Start(t *testing.T) {
    Start("127.0.0.1:9090", "127.0.0.1:9191")
}
