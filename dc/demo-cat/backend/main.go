package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zaynjarvis/fyp/dc/sdk"
)

func main() {
	s := sdk.New("agent-cat:7000", "cat")
	go func() {
		if err := s.Connect(context.Background()); err != nil {
			log.Printf("example failed, err: %v\n", err)
		}
	}()

	images := make([][]byte, 0, 5)
	for i := range make([]string, 5) {
		img, err := ioutil.ReadFile(fmt.Sprintf("./assets/%d.jpg", i+1))
		if err != nil {
			panic(err)
		}
		images = append(images, img)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req struct {
			Pic   int     `json:"pic"`
			Score float64 `json:"score"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("%#v\n", req)

		s.Image(fmt.Sprintf("%s.jpg", time.Now().Format("15-04-05")), images[req.Pic], req)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "{}")
	})

	log.Println("server listening on 9090")
	log.Fatal(http.ListenAndServe(":9090", cors.Default().Handler(mux)))
}
