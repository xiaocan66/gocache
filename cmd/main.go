package main

import (
	"flag"
	"fmt"
	"gocache"
	"log"
	"net/http"
)

var db = map[string]string{
	"laowang": "122333",
	"li":      "lizican123@gmail.com",
	"jack":    "hello worlds",
	"tom":     "woshi nibaba ",
}

func createGroup() *gocache.Group {
	return gocache.NewGroup("user", 2<<10, gocache.GetterFunc(func(key string) ([]byte, error) {
		// 从数据库中读取数据
		log.Println("[slowDb] search key ", key)

		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s no exist", key)
	}))

}

func startCacheServer(addr string, addrs []string, gee *gocache.Group) {
	peers := gocache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)

	log.Fatal(http.ListenAndServe(addr[7:], peers))

}

func startApiServer(apiAddr string, gcc *gocache.Group) {
	http.Handle("/api", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("key")
		view, err := gcc.Get(key)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/octet-stream")
		_, err = writer.Write(view.ByteSlice())
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gocache server port")
	flag.BoolVar(&api, "api", false, "start a api server")
	flag.Parse()
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
		8004: "http://localhost:8004",
		8005: "http://localhost:8005",
	}
	gcc := createGroup()
	if api {
		go startApiServer(apiAddr, gcc)
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	startCacheServer(addrMap[port], addrs, gcc)
}
