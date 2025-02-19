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

var FILE_NAME string
