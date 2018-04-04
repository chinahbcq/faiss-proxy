package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "github.com/chinahbcq/faiss-proxy/proto"
)

var (
	faissServer = flag.String("faiss_endpoint", "10.1.34.159:3838", "endpoint of Faiss gRPC Service")
	port        = flag.Int("http_port", 3839, "Port of http server")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	pat, _ := runtime.NewPattern(
		1,
		[]int{
			int(utilities.OpLitPush), 0,
			int(utilities.OpLitPush), 1,
			int(utilities.OpLitPush), 2,
		},
		[]string{"faiss", "1.0", "swagger.json"},
		"",
	)
	mux.Handle("GET", pat, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		io.Copy(w, strings.NewReader(gw.Swagger))
	})

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterFaissServiceHandlerFromEndpoint(ctx, mux, *faissServer, opts)
	if err != nil {
		return err
	}

	portStr := fmt.Sprintf("%d", *port)
	log.Print("Faiss gRPC Server gateway start at port " + portStr + "...")
	err = http.ListenAndServe(":"+portStr, mux)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
