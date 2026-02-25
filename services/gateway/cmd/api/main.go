package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

const storage string = "http://storage:4000"

func main() {
	storageUrl := flag.String("storage", storage, "Address of the storage service")
	flag.Parse()

	server := &http.Server{
		Addr: ":4000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp, err := http.Get(*storageUrl)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte(fmt.Sprintf("success: %s", string(body))))
		}),
	}

	fmt.Println("gateway started")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
