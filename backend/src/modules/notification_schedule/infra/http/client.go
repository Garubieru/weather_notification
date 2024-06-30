package http

type Client interface {
	Request(request Request)
}

type Request struct {
	Url    string
	Method string
}
