package main

import "net/http"

var COOKIES http.Cookie = http.Cookie{
	Name:     "session",
	Quoted:   false,
	MaxAge:   0,
	Secure:   true,
	HttpOnly: true,
	SameSite: http.SameSiteDefaultMode,
}

const RESET string = "\033[0m"
const BLUE string = "\033[34m"
