package server

var VERSION string

func (srv *Server) Version() APIVersion {
	return APIVersion{
		Version: VERSION,
	}
}

func NewServer() (*Server, error) {
	srv := &Server{}
	return srv, nil
}

type Server struct{}
