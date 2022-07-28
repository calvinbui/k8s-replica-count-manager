package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes"
)

func (s *server) StartHTTPHealthCheck(ctx context.Context) error {
	// endpoint and handler for the http health check
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := testK8sConnection(s.services.K8s, r); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "not ok!")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok!")
	})

	// start http server
	if err := http.ListenAndServe(s.config.HttpAddress, nil); err != nil {
		return err
	}

	return nil
}

func testK8sConnection(client *kubernetes.Clientset, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// contact the Kubernetes' health endpoint to check connection
	res, err := client.Discovery().RESTClient().Get().AbsPath("/healthz").DoRaw(ctx)
	if err != nil {
		return err
	}

	if string(res) != "ok" {
		return err
	}

	return nil
}
