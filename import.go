package main

// Implements importing history from a Slack export from:
// https://my.slack.com/services/export

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Message struct {
	Type string
	User string
	Text string
	TS   string
}

var (
	messageRegex *regexp.Regexp
)

func init() {
	messageRegex = regexp.MustCompile(`<([^>]+)>`)
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
		log.Print("Import complete")
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

func parseText(text string) string {
	matches := messageRegex.FindAllStringSubmatch(text, -1)
	for _, matches2 := range matches {

		if strings.HasPrefix(matches2[1], "http") || strings.HasPrefix(matches2[1], "mailto") {
			strings.Replace(text, matches2[0], "", -1)

		} else if strings.HasPrefix(matches2[1], "@U") {
			parts := strings.SplitN(matches2[1], "|", 2)

			if len(parts) == 2 {
				text = strings.Replace(text, matches2[0], parts[1], -1)
			} else {
				text = strings.Replace(text, matches2[0], "", -1)
			}

		} else if strings.HasPrefix(matches2[1], "@") {
			text = strings.Replace(text, matches2[0], matches2[1], -1)

		} else if strings.HasPrefix(matches2[1], "#") {
			parts := strings.SplitN(matches2[1], "|", 2)

			if len(parts) == 2 {
				text = strings.Replace(text, matches2[0], parts[1], -1)
			} else {
				text = strings.Replace(text, matches2[0], "", -1)
			}

		}

		text = strings.TrimSpace(text)

		text = strings.Replace(text, "&lt;", "<", -1)
		text = strings.Replace(text, "&gt;", ">", -1)
		text = strings.Replace(text, "&amp;", "&", -1)
	}

	return text
}
