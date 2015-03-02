package main

// Implements importing history from a Slack export from:
// https://my.slack.com/services/export

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Message struct {
	Type string
	User string
	Text string
	TS   string
}

// Start an import from a slack export directory
// Does some basic error checking then imports the data in the background
func StartImport(dir *string, channel *string) (err error) {
	log.Printf("Starting import from %s", *dir)

	// Does this directory exist? Get its contents
	contents, err := ioutil.ReadDir(*dir)
	if err != nil {
		return err
	}

	// Is there a users.json file?
	_, err = ioutil.ReadFile(*dir + "/users.json")
	if err != nil {
		return err
	}

	// Looks good, import each directory/channel in a goroutine
	go func() {
		for _, file := range contents {
			if file.IsDir() {
				if *channel == "" || *channel == file.Name() {
					ImportDir(*dir + "/" + file.Name())
				}
			}
		}

		// Write the state file
		err = markovChain.Save(stateFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Import complete. %d suffixes", len(markovChain.Chain))
	}()

	return nil
}

// Handles the import of a channel/directory
func ImportDir(dir string) {
	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range contents {
		if !file.IsDir() {
			contents, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			// parse into json
			var messages []Message
			err = json.Unmarshal(contents, &messages)
			if err != nil {
				log.Fatal(err)
			}

			for _, message := range messages {
				text := message.Text
				text = parseText(text)
				if text != "" {
					markovChain.Write(text)
				}
			}
		}
	}
}
