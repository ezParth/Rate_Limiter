package reverseproxy

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	config "rl/configuration"
	"sync"
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

func RunProxyServer(proxy string) []byte {
	// creating a request
	req, err := http.NewRequest("GET", proxy, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// making the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return body
}

func Start(w http.ResponseWriter, r *http.Request) {
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

	proxyServer, err := getProxyServer(yml, node)
	if err != nil {
		panic(err)
	}

	fmt.Println("Proxy Server: ", proxyServer)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		reply := RunProxyServer(proxyServer)
		if reply == nil {
			_, err := w.Write([]byte("no reply"))
			if err != nil {
				http.Error(w, "failed to write response", http.StatusInternalServerError)
				return
			}
		} else {
			_, err := w.Write(reply)
			if err != nil {
				http.Error(w, "failed to write response", http.StatusInternalServerError)
				return
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
