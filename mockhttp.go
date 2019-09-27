package mockhttp

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)


type MockHttpServer struct {
	server http.Server
	handler func(w http.ResponseWriter, r *http.Request)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func randomAddr() string {
	for {
		port := 1024 + rand.Int31n(65535 - 1024)
		addr := fmt.Sprintf("localhost:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return addr
		}
	}
}

func MakeMockHttpServer() *MockHttpServer {
	var mockHttp MockHttpServer

	// Just print 'ok'
	mockHttp.AssignDefaultHandler()

	// Create a router that listens to all routes,
	// handled by MockHttpServer.handler
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockHttp.handler(w, r)
	})

	// Create the server.
	mockHttp.server = http.Server{
		Addr:              randomAddr(),
		Handler:           r,
	}

	// BUG: Chance of bind() on Addr by other processes.
	go func() {
		mockHttp.server.ListenAndServe()
	}()

	for i := 0; i < 11; i++ {
		_, err := http.Get(mockHttp.GetBaseUri())

		if err == nil {
			break
		} else {
			time.Sleep(time.Millisecond * 100)
		}

		if i == 10 {
			log.Fatalf("Unable to bind port %s: %s", mockHttp.GetBaseUri(), err.Error())
		}
	}

	return &mockHttp
}

func (m *MockHttpServer) GetBaseUri() string {
	return "http://" + m.server.Addr + "/"
}

func (m *MockHttpServer) AssignDefaultHandler() {
	m.handler = defaultHandler
}

func (m *MockHttpServer) AssignHandler(f func(w http.ResponseWriter, r *http.Request)) {
	m.handler = f
}

func (m *MockHttpServer) Close() {
	err := m.server.Close()
	if err != nil {
		log.Fatalf("Error during close %s", err)
	}
}

