package main

import (
	"log"
	"regexp"
)

var forbidden_chars *regexp.Regexp

func init() {
	regx, err := regexp.Compile(`\W`) // Nb this is 'not word' (ie not one of [A-Za-z0-9_])
	if err != nil {
		// If we can't remove chars don't allow start
		log.Fatalln(err)
	}
	forbidden_chars = regx
}

// Scrub the hell out of an incoming string, return resulting string.
func scrub(s string) string {
	return forbidden_chars.ReplaceAllString(s, "")
}
