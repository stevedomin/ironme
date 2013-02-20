package main

import (
	"flag"
	"fmt"
	"github.com/ngmoco/falcore"
	"github.com/ngmoco/falcore/compression"
	"github.com/stevedomin/ironme/coffee_file"
	"github.com/stevedomin/ironme/sass_file"
	"github.com/stevedomin/ironme/static_file"
	"net/http"
)

var (
	port                 = flag.Int("p", 4000, "Server port")
	singlePageApp        = flag.Bool("-s", true, "Server port")
	path          string = "."
)

func main() {
	// Parse command line options
	flag.Parse()

	if flag.Arg(0) != "" {
		path = flag.Arg(0)
	}

	// setup pipeline
	pipeline := falcore.NewPipeline()

	// upstream filters
	// Serve index.html for root requests
	pipeline.Upstream.PushBack(falcore.NewRequestFilter(func(req *falcore.Request) *http.Response {
		if req.HttpRequest.URL.Path == "/" {
			req.HttpRequest.URL.Path = "/index.html"
		}
		return nil
	}))

	// Serve files
	pipeline.Upstream.PushBack(&static_file.Filter{
		BasePath: path,
	})

	// Serve SASS files
	pipeline.Upstream.PushBack(&sass_file.Filter{
		BasePath: path,
	})

	// Serve Coffee files
	pipeline.Upstream.PushBack(&coffee_file.Filter{
		BasePath: path,
	})

	if *singlePageApp {
		// Serve index.html no matter what
		pipeline.Upstream.PushBack(falcore.NewRequestFilter(func(req *falcore.Request) *http.Response {
			req.HttpRequest.URL.Path = "/index.html"
			return nil
		}))

		// Rewrite me !
		pipeline.Upstream.PushBack(&static_file.Filter{
			BasePath: path,
		})
	}

	// downstream
	pipeline.Downstream.PushBack(compression.NewFilter(nil))

	// setup server
	server := falcore.NewServer(*port, pipeline)

	// start the server
	// this is normally blocking forever unless you send lifecycle commands 
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}
}
