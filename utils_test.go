package main

import (
	"testing"
)

func TestParseText(t *testing.T) {
	var text string

	text = parseText("It’s a really really good book")
	if text != "It’s a really really good book" {
		t.Error("Plain text parsing failed")
	}

	text = parseText("<https://twitter.com/mikeash/status/645014517005570048>")
	if text != "" {
		t.Error("URL stripping failed")
	}

	text = parseText("&gt;&gt;&gt; \"And our cardboard cutout of Niall from One Direction had been moved — someone had taken him out of the living room and put him in the garage.\"")
	if text != ">>> \"And our cardboard cutout of Niall from One Direction had been moved — someone had taken him out of the living room and put him in the garage.\"" {
		t.Error("Entity replacement failed")
	}

	text = parseText("That should be posted in <#C05476P6Z>")
	if text != "That should be posted in" {
		t.Error("Bare channels failed")
	}

	text = parseText("That should be posted in <#C05476P6Z|random>")
	if text != "That should be posted in #random" {
		t.Error("Channel with name failed")
	}

	text = parseText("That sounds like something for <@U06CTQTRU|myles>")
	if text != "That sounds like something for @myles" {
		t.Error("User with name failed")
	}

	text = parseText("That sounds like something for <@myles>")
	if text != "That sounds like something for @myles" {
		t.Error("User with with no id failed")
	}
}
