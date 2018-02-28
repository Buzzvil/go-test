# go-test

This is go library project for test cases.

## mock server

You can easily make a mock response for a http request by using this.

```go

//Define httpClient you want to add a patch.
httpClient := http.DefaultClient
//Define a target server with domain and add ResponseHandlers. ResponseHandler contains a mock response and a request should be handled.
testServer := mock.NewTargetServer("google.com").AddResponseHandler(&mock.ResponseHandler{
	WriteToBody: func() []byte {
		return []byte(testBody)
	},
	Path:       "/",
	Method:     http.MethodGet,
	StatusCode: 400,
})
//After PatchClient, every http request using the httpClient will handle the mock request.
clientPatcher := mock.PatchClient(httpClient, testServer)
//RemovePatch will remove all mock handlers from the httpClient.
defer clientPatcher.RemovePatch()
```
