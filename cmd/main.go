package main

import (
	"context"
	"flag"
	"github.com/joker-bai/validate-namespace/http"
	log "k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

var (
	tlscert, tlskey, port string
)

func main() {
	flag.StringVar(&tlscert, "tlscert", "/etc/certs/cert.pem", "Path to the TLS certificate")
	flag.StringVar(&tlskey, "tlskey", "/etc/certs/key.pem", "Path to the TLS key")
	flag.StringVar(&port, "port", "8443", "The port to listen")
	flag.Parse()

	server := http.NewServer(port)
	go func() {
		if err := server.ListenAndServeTLS(tlscert, tlskey); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	log.Infof("Server running in port: %s", port)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info("Shutdown gracefully...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
