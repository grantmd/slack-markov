package main

// Implements importing history from a Slack export from:
// https://my.slack.com/services/export

import (
	"io/ioutil"
	"log"
	"time"
)

func StartImport(dir *string) (err error) {
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

	// Looks good, import each directory in a goroutine
	go func() {
		for _, file := range contents {
			time.Sleep(5 * time.Second)
			if file.IsDir() {
				log.Printf("Directory: %s", file.Name())
				ImportDir(*dir + "/" + file.Name())
			} else {
				log.Printf("File: %s", file.Name())
			}
		}
	}()

	return nil
}

func ImportDir(dir string) {

}
