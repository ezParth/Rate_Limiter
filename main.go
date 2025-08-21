package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	config "rl/configuration"
	"syscall"
	"time"
)

type Upstream struct {
	url  string
	node []string
}

func getUpstreams(file *config.YAMLCONFIG, url string) (*Upstream, error) {
	foundUrl := false
	upstream := &Upstream{url: url}
	for _, rules := range file.Server.Rules {
		if rules.Path == url {
			foundUrl = true
			upstream.node = rules.Upstream
		}
	}

	if !foundUrl {
		return nil, fmt.Errorf("cannot find Url '%s' in configurations", url)
	}

	return upstream, nil
}

func getUpstream(upstream *Upstream) (string, error) {
	size := len(upstream.node)
	if size == 0 {
		return "", fmt.Errorf("Upstream is empty, no node assigned")
	}
	return upstream.node[rand.Intn(size)], nil
}

func getProxyServer(file *config.YAMLCONFIG, node string) (string, error) {
	var server []string
	for _, servers := range file.Server.Upstream {
		if servers.ID == node {
			server = append(server, servers.Server)
		}
	}
	if len(server) == 0 {
		return "", fmt.Errorf("cannot Find Node %s", node)
	}

	return server[rand.Intn(len(server))], nil
}

func start(w http.ResponseWriter, r *http.Request) {
	yml, err := config.LoadYAML("config.yml")
	if err != nil {
		panic(err)
	}

	url := r.URL.String()

	upstream, err := getUpstreams(yml, url)
	if err != nil {
		panic(err)
	}

	node, err := getUpstream(upstream)

	if err != nil {
		panic(err)
	}

	fmt.Println("Node: ", node)

	proxyServer, err := getProxyServer(yml, node)
	if err != nil {
		panic(err)
	}

	fmt.Println("Proxy Server: ", proxyServer)

	w.Write([]byte("Hello, from server"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", start)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Printf("Server started in port %s\n", server.Addr[1:])

		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
