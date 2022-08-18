package main

import "testing"

func TestParseLink(t *testing.T) {
	res := Reference{}

	link := "++ https://aaa.com"
	logger.Info("expected it", "star", 2, "link", "https://aaa.com")
	res.parseLink(link)
	logger.Info("parse result", "star", res.Good, "link", res.Link)

	link = "http://bbb.com/aaa"
	logger.Info("expected it", "star", 0, "link", link)
	res.parseLink(link)
	logger.Info("parse result", "star", res.Good, "link", res.Link)
}
