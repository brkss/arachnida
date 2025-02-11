package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brkss/arachnida/spider/internal/infrastructure"
	"github.com/brkss/arachnida/spider/internal/usecase"
)

func main() {

	var (
		recursive bool
		depth     int
		path      string
	)

	flag.BoolVar(&recursive, "r", false, "recursive")
	flag.IntVar(&depth, "d", 0, "depth")
	flag.StringVar(&path, "p", "", "path")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: spider -r -d <depth> -p <path>")
		os.Exit(1)
	}

	startURL := args[0]

	// validate depth 
	finalDepth, err := usecase.ValidateDepth(depth)
	if err != nil {
		log.Fatalf("failed to validate depth: %v", err)
	}


	client := &http.Client{}
	
	spiderService := infrastructure.NewSpiderService(client)
	fileSaver := infrastructure.NewFileSaver(client)
	spiderUsecase := usecase.NewSpiderUsecase(spiderService, fileSaver)

	err = spiderUsecase.DownloadImages(startURL, finalDepth, 1, path)
	if err != nil {
		log.Fatalf("failed to download images: %v", err)
	}


	fmt.Println("Downloaded images successfully");
	
	
}