package main

import (
	"contactsapi/api"
	"net/http"
)

func main() {
	api.Init()
	http.ListenAndServe(":8080", api.Handlers())
}
