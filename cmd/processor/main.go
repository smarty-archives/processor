package main

import (
	"log"
	"net/http"
	"os"

	"github.com/smartystreets/processor"
)

func main() {
	client := processor.NewAuthenticationClient(
		http.DefaultClient, "https", "us-street.api.smartystreets.com",
		"e67804f5-839b-7f81-0ade-63faa39e7389", "IAAUwgcmDrM4FpLm94eH")

	pipeline := processor.NewPipeline(os.Stdin, os.Stdout, client, 8)

	if err := pipeline.Process(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
