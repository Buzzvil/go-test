package mock_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Buzzvil/go-test/mock"
)

func TestMockRequest(t *testing.T) {
	testBody := "Hello world"

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

	res := validateAndGetResponse(t, httpClient, 400)

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	} else if testBody != string(bodyData) {
		t.Fatalf("TestMockRequest() - test: %s, body: %s", testBody, bodyData)
	}
	//RemovePatch will remove all mock handlers from the httpClient.
	clientPatcher.RemovePatch()

	validateAndGetResponse(t, httpClient, 200)
}

func validateAndGetResponse(t *testing.T, httpClient *http.Client, statusCode int) *http.Response {
	req, err := http.NewRequest("GET", "http://google.com/", nil)

	if err != nil {
		panic(err)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		panic(err)
	}

	if err != nil || res.StatusCode != statusCode {
		t.Fatal(res.StatusCode, err)
	}
	return res
}
