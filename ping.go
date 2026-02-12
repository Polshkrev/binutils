package binutils

import (
	"net/http"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections/safe"
	"github.com/Polshkrev/otvet"
)

// Concurrently send an http request to a given url.
func getResponse(url string, response chan<- *http.Response, errorChannel chan<- error) {
	var result *http.Response
	var responseError error
	result, responseError = http.Get(url)
	response <- result
	errorChannel <- responseError
	defer close(response)
	defer close(errorChannel)
}

// Ping a given url.
// Returns an http response.
// If the http request fails, an [gopolutils.Exception] is returned with a nil data pointer.
func ping(url string) (*http.Response, *gopolutils.Exception) {
	var responseChannel chan *http.Response = make(chan *http.Response, 1)
	var errorChannel chan error = make(chan error, 1)
	go getResponse(url, responseChannel, errorChannel)
	var response *http.Response = <-responseChannel
	var responseError error = <-errorChannel
	if responseError != nil {
		return nil, gopolutils.NewException(responseError.Error())
	}
	return response, nil
}

// Ping the status code of given urls.
// Returns a collection of http status codes.
// If the http request fails, an [gopolutils.Exception] is returned with a nil data pointer.
func Ping(urls ...string) (collections.Collection[otvet.StatusCode], *gopolutils.Exception) {
	var codes safe.Collection[otvet.StatusCode] = safe.NewArray[otvet.StatusCode]()
	var i int
	for i = range urls {
		var url string = urls[i] // Idiomatically more performant than direct access.
		var response *http.Response
		var except *gopolutils.Exception
		response, except = ping(url)
		if except != nil {
			return nil, except
		}
		codes.Append(otvet.StatusCode(response.StatusCode)) // ! If the status code is not in the enum list, this will panic.
	}
	return codes, nil
}
