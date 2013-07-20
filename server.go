package gsix

import (
	"net/http"
	"log"
	"fmt"
)

type Server struct {}

func NewServer() (s *Server) {
	server := new(Server)

	return server
}

func (s *Server) Listen(port uint) {
	connectStr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(connectStr, nil)	)
}
