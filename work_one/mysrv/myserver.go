package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func init() {
	os.Setenv("VERSION", "v0.0.1")
}

type Output struct {
	Code int
	Msg  string
}

func srvApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)

	return http.ListenAndServe("0.0.0.0:18080", mux)
}

func index(resp http.ResponseWriter, req *http.Request) {
	version := os.Getenv("VERSION")
	resp.Header().Set("VERSION", version)
	resp.Header().Set("content-type", "application/json")

	for k, v := range req.Header {
		for _, vv := range v {
			resp.Header().Set(k, vv)
		}
	}

	output := Output{
		Code: 200,
		Msg:  "ok",
	}

	b, _ := json.Marshal(output)
	resp.Write(b)

	clientIP, err := getClientIP(req)
	if err != nil {
		log.Printf("Faild to get client IP: %+v\n", err)
	}

	log.Printf("Success! Client IP: %v\n", clientIP)
}

func getClientIP(req *http.Request) (string, error) {
	var err error
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip, nil
	}

	ip = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip, nil
	}

	if ip, _, err = net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip, nil
	}

	return "", err
}

func healthz(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Working ...")
}

// Signal ...
func Signal() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Program Exiting...", sig)
				done <- true
			default:
				log.Println("Other signal", sig)
			}
		}
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	log.Println("Awaiting signal")
	<-done
	log.Println("Exited")
}

func main() {
	go func() {
		err := srvApp()
		if err != nil {
			log.Fatalf("start http server failed, error: %s\n", err.Error())
		}
	}()
	Signal()
}
