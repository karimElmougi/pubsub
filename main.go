package main

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/karimElmougi/pubsub/internal/pubsub"
)

func main() {
	var port int
	if len(os.Args) < 2 {
		port = 0
	} else {
		var err error
		port, err = strconv.Atoi(os.Args[1])
		if err != nil || port < 0 || port > math.MaxUint16 {
			fmt.Println("invalid port: ", os.Args[1])
			os.Exit(1)
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subgroup := pubsub.NewSubscriberGroup()
	handler := pubsub.NewServer(subgroup)

	fmt.Println("Starting server on port: ", listener.Addr().(*net.TCPAddr).Port)
	err = http.Serve(listener, handler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
