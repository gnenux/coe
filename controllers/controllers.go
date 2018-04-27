package controllers

var URL map[string]string

func init() {
	URL = map[string]string{}
	URL["SearchURL"] = "/search"
	URL["Companies"] = "/companies"
}
