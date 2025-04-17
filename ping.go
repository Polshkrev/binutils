package binutils

import (
	"net/http"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
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
// If the http request fails, an Exception is returned with a nil data pointer.
func ping(url string) (*http.Response, *gopolutils.Exception) {
	var responseChannel chan *http.Response
	var errorChannel chan error
	go getResponse(url, responseChannel, errorChannel)
	var response *http.Response = <-responseChannel
	var responseError error = <-errorChannel
	if responseError != nil {
		return nil, gopolutils.NewException(responseError.Error())
	}
	return response, nil
}

// Ping the status code of given urls.
// Return a collection of http status codes.
// If the http request fails, an Exception is returned with a nil data pointer.
func Ping(urls ...string) (collections.Collection[uint16], *gopolutils.Exception) {
	var url string
	var codes collections.Collection[uint16] = collections.NewArray[uint16]()
	for _, url = range urls {
		var response *http.Response
		var except *gopolutils.Exception
		response, except = ping(url)
		if except != nil {
			return nil, except
		}
		codes.Append(uint16(response.StatusCode))
	}
	return codes, nil
}
