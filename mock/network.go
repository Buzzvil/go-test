package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ClientPatcher struct {
	OriginalTransport http.RoundTripper
	MockTransport     http.RoundTripper
	httpClient        *http.Client
}

func PatchClient(httpClient *http.Client, targetServers ...*targetServer) *ClientPatcher {

	t := &Transport{
		OriginalTransport: httpClient.Transport,
		TargetServers:     targetServers,
	}

	if t.OriginalTransport == nil {
		t.OriginalTransport = &http.Transport{}
	}

	patcher := &ClientPatcher{
		OriginalTransport: httpClient.Transport,
		MockTransport:     t,
		httpClient:        httpClient,
	}

	httpClient.Transport = patcher.MockTransport

	return patcher
}

func NewTargetServer(host string) *targetServer {
	if host == "" {
		panic("Host should be defined.")
	}

	return &targetServer{
		Host: host,
	}
}

func (s *targetServer) AddResponseHandler(r *ResponseHandler) *targetServer {
	if r.WriteToBody == nil {
		panic("WriteToBody of resHandler should be defined.")
	}

	if r.StatusCode == 0 {
		r.StatusCode = 200
	}

	if r.Method == "" {
		panic("Method of resHandler should be defined.")
	}

	s.ResponseHandlers = append(s.ResponseHandlers, r)
	return s
}

func (httpReq *ClientPatcher) RemovePatch() {
	httpReq.httpClient.Transport = httpReq.OriginalTransport
}

type targetServer struct {
	Host             string
	ResponseHandlers []*ResponseHandler
}

type ResponseHandler struct {
	WriteToBody func() []byte
	StatusCode  int
	Path        string
	Method      string
}

type RequestHandler interface {
	GetResponse() []byte
}

type Transport struct {
	TargetServers     []*targetServer
	OriginalTransport http.RoundTripper
}

func (m *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, targetServer := range m.TargetServers {
		if targetServer.Host == req.Host {
			for _, resHandler := range targetServer.ResponseHandlers {
				if resHandler.Path == req.URL.Path && resHandler.Method == req.Method {
					// mock response object
					r := http.Response{
						StatusCode: resHandler.StatusCode,
						Proto:      "HTTP/1.0",
						ProtoMajor: 1,
						ProtoMinor: 0,
					}

					buf := bytes.NewBuffer(resHandler.WriteToBody())
					r.Body = ioutil.NopCloser(buf)

					return &r, nil
				}
			}
		}
	}
	return m.OriginalTransport.RoundTrip(req)
}
