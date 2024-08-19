package ports

import "net/http"

type ServerPort interface {
	PrepareServer() (*http.Server, error)
}
