slack-markov
=======

A markov chain bot for [Slack](https://slack.com) written in Go (golang)

[![Build Status](https://travis-ci.org/grantmd/slack-markov.png)](https://travis-ci.org/grantmd/slack-markov)

Usage
-----

* Build the code with `go build`

* Start the bot with `./slack-markov` on an internet-accessible server. (Check the output of `./slack-markov -h` for configuration options)

* Configure an [Outgoing Webhook](https://my.slack.com/services/new/outgoing-webhook) in your Slack and point it to the place where your bot is running. For example: `http://example.com:8000/`

* The bot will listen to incoming requests and add the messages to the chain. At a configurable response rate, it will also build and send a response to the channel

Tips
----

* Export your team's data from https://my.slack.com/services/export, and use that to seed the chain. See the `-importDir` and `-importChan` options.
* Keep your bot scoped to one channel.
* Adjust the `-responseChance` option based on your team's activity and how much you want to annoy your co-workers
* Look into the `-twitter*` options if you enjoy Twitter and libel