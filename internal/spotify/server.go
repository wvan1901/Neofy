package spotify

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type loginServer struct {
	server   *http.Server
	codeChan chan string
	done     chan bool
}

func CreateLoginServer(codeChan chan string) *loginServer {
	doneChan := make(chan bool)
	srv := &http.Server{Addr: ":8090"}
	http.HandleFunc("/callback", callbackHandler(codeChan, doneChan))
	return &loginServer{
		server:   srv,
		codeChan: codeChan,
		done:     doneChan,
	}
}

func callbackHandler(codeChan chan string, done chan bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			codeChan <- code
			fmt.Println("Found Code")
			// Send to chan that we are done (so shut server down)
			defer func() {
				done <- true
			}()
			resBuf := bytes.NewBufferString("You are authenticated, close tab and return to cli")
			w.Write(resBuf.Bytes())
		} else {
			resBuf := bytes.NewBufferString("{'Status':'error'}")
			w.Write(resBuf.Bytes())
		}
	}
}

func (s *loginServer) startHttp() {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
	fmt.Println("Terminated http server")
}

func (s *loginServer) endHttp() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("end server: %v", err)
	}
}

func (s *loginServer) RunServer() {
	go s.startHttp()
	go s.isDone()
}

func (s *loginServer) isDone() {
	select {
	case <-s.done:
		// Give some time for server to respond then shutdown
		time.Sleep(5 * time.Second)
		s.endHttp()
	}
}
