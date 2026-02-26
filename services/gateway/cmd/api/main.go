package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	storageUrl := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_PORT"))
	selfUrl := fmt.Sprintf(":%s", os.Getenv("GATEWAY_PORT"))

	server := &http.Server{
		Addr: selfUrl,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp, err := http.Get(storageUrl)
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

	log.Println("gateway started on ", selfUrl)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
