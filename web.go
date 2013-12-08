package main

// Slack outgoing webhooks are handled here. Requests come in and are run through
// the markov chain to generate a response, which is sent back to Slack.
//
// Create an outgoing webhook in your Slack here:
// https://my.slack.com/services/new/outgoing-webhook

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type WebhookResponse struct {
	Text string `json:"text"`
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Handling incoming request")
		markovChain.Write(r.PostFormValue("text"))
		go func() {
			markovChain.Save(stateFile)
		}()

		var response WebhookResponse
		response.Text = markovChain.Generate(numWords)

		b, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(b)
	})
}

func StartServer(port int) {
	log.Printf("Starting HTTP server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
