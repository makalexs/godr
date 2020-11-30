package httpserver

import (
	"github.com/makalexs/godr/httpserver/api"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  {
	return c.in.Read(p)
}

func (c *HttpConn) Write(d []byte) (n int, err error) {
	return c.out.Write(d)
}

func (c *HttpConn) Close() error {
	return nil
}

func StartServer () (int,error) {

	server := rpc.NewServer()
	err := server.Register(&api.Common{})
	if err != nil {
		log.Fatalln(err)
	}
	//port := os.Getenv("DR_COMMON_API_PORT")
	listener, err := net.Listen("tcp", ":" + "12345")

	if err != nil {
		//panic(err)
		return 0, err
	}

	defer listener.Close()

	err = http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rpc/common" {
			serverCodec := jsonrpc2.NewServerCodec(&HttpConn{in: r.Body, out: w}, server)

			w.Header().Set( "Content-type", "application/json" )
			w.WriteHeader(200)

			if err = server.ServeRequest(serverCodec) ; err != nil {
				http.Error(w, "Error while serving JSON request", 500)
				return
			}
		} else {
			http.Error(w, "Unknown request", 404)
		}
	} ) )
	if err != nil {
		log.Fatalln(err)
	}
	return 1, nil

}
