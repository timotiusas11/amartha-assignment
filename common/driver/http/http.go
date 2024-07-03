package http

type HTTPInterface interface {
	Post() error
}

type HTTP struct {
	// Client
}

func New() HTTP {
	return HTTP{}
}

func (h HTTP) Post() error {
	// Call downstream using POST method
	return nil
}
