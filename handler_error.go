package main

import "net/http"

func handlerError(w http.ResponseWriter, r *http.Request) {
	repondWithError(w, 400, "Something went wrong.	")
}
