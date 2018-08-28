package httplib

import (
	"bytes"
	"crypto/tls"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

)

// DefaultClient is used if no custom HTTP client is defined.
var DefaultClient *Client = NewClient()

// Library level defaults.
var (
	defaultClientTimeout         = 180 * time.Second
	defaultDialTimeout           = 30 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultMaxIdleConns          = 100
)

// SendRaw makes a request based on http.Request using DefaultClient.
func SendRaw(req *http.Request) (*http.Response, error) {
	return DefaultClient.SendRaw(req)
}

// Send makes request based on Request using DefaultClient.
func Send(req *Request) (*Response, error) {
	return DefaultClient.Send(req)
}

// Get makes Get request using DefaultClient.
func Get(url string) (*Response, error) {
	return DefaultClient.Get(url)
}

// Post makes Post request using DefaultClient.
func Post(url, contentType string, body []byte) (*Response, error) {
	return DefaultClient.Post(url, contentType, body)
}

// DownloadFile performs DownloadFile using DefaultClient.
func DownloadFile(url, filename string) (int, error) {
	return DefaultClient.DownloadFile(url, filename)
}

// Request is a caller facing type to create a request.
type Request struct {
	Method  string
	BaseURL string
	Body    []byte
	Header  http.Header
	Param   url.Values
}

// NewRequest returns Request ready for the usage.
func NewRequest(m string, u string, b []byte) *Request {
	return &Request{
		Method:  m,
		BaseURL: u,
		Body:    b,
		Header:  make(http.Header),
		Param:   make(url.Values),
	}
}

// AddParamToURL returns the url with params.
func (r *Request) AddParamToURL() string {
	return tools.JoinStrings(r.BaseURL, "?", r.Param.Encode())
}

// SetContentTypeJSON sets header Content-Type to "application/json".
func (r *Request) SetContentTypeJSON() {
	if r.Header == nil {
		r.Header = make(http.Header)
	}

	r.Header.Set("Content-Type", "application/json")
}

// SetContentType sets header Content-Type to a given value.
func (r *Request) SetContentType(value string) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}

	r.Header.Set("Content-Type", value)
}

// SetAuthorization sets header Authorization to a given value.
func (r *Request) SetAuthorization(value string) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}

	r.Header.Set("Authorization", value)
}

// Response holds the response from the call.
type Response struct {
	StatusCode int
	Status     string
	Body       []byte
	Headers    http.Header
}

// Client is a properly set up http client with methods.
type Client struct {
	HTTPClient http.Client
}

// NewClient returns Client with customized default transport and timeouts.
//
// More details here https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
func NewClient() *Client {
	return newClientWithSettings(
		defaultDialTimeout, defaultKeepAlive,
		defaultIdleConnTimeout, defaultTLSHandshakeTimeout, defaultExpectContinueTimeout,
		defaultClientTimeout,
		true, true,
	)
}

// NewClientWithSettings allows to create a Client with custom settings.
//
// WARNING! You MUST totally understand what are you doing by passing incorrect values.
func NewClientWithSettings(
	// Dialer timouts.
	dialTimeout, keepAlive time.Duration,
	// Transport timouts.
	idleTimeout, tlsHandshakeTimeout, continueTimeout time.Duration,
	// Client timeouts.
	clientTimeout time.Duration,
	// Transport TLS settings.
	skipTLSVerify, customCACerts bool,
) *Client {
	return newClientWithSettings(
		dialTimeout, keepAlive,
		idleTimeout, tlsHandshakeTimeout, continueTimeout,
		clientTimeout,
		skipTLSVerify, customCACerts,
	)
}

// newClientWithSettings is a private constructor used internally.
func newClientWithSettings(
	// Dialer timouts.
	dialTimeout, keepAlive time.Duration,
	// Transport timouts.
	idleTimeout, tlsHandshakeTimeout, continueTimeout time.Duration,
	// Client timeouts.
	clientTimeout time.Duration,
	// Transport TLS settings.
	skipTLSVerify, customCACerts bool,
) *Client {
	// Create a custom tls config.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipTLSVerify,
	}

	// Create a pool with root CA.
	if customCACerts {
		certPool, err := CACerts()
		if err == nil {
			tlsConfig.RootCAs = certPool
		}
	}

	// Create a custom transport.
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: keepAlive,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       idleTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: continueTimeout,
		TLSClientConfig:       tlsConfig,
	}

	// Finally, create a custom http client.
	client := &Client{
		HTTPClient: http.Client{
			Timeout:   clientTimeout,
			Transport: transport,
		},
	}

	return client
}

// buildRequest creates the http.Request from Request.
func (c *Client) buildRequest(req *Request) (*http.Request, error) {
	if len(req.Param) != 0 {
		req.BaseURL = req.AddParamToURL()
	}

	r, err := http.NewRequest(req.Method, req.BaseURL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	if len(req.Header) != 0 {
		r.Header = req.Header
	}

	return r, nil
}

// SendRaw performs a request based on http.Request.
func (c *Client) SendRaw(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// buildResponse builds Response from http.Response.
//
// It takes care of closing body as well.
func (c *Client) buildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := &Response{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Body:       body,
		Headers:    res.Header,
	}

	return response, nil
}

// Send makes request.
func (c *Client) Send(request *Request) (*Response, error) {
	// Build the HTTP request object.
	req, err := c.buildRequest(request)
	if err != nil {
		return nil, err
	}

	// Build the HTTP client and make the request.
	res, err := c.SendRaw(req)
	if err != nil {
		return nil, err
	}

	// Build Response object.
	response, err := c.buildResponse(res)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get makes GET request.
func (c *Client) Get(url string) (*Response, error) {
	req := NewRequest("GET", url, nil)

	return c.Send(req)
}

// Post makes Post request.
func (c *Client) Post(url, contentType string, body []byte) (*Response, error) {
	req := NewRequest("POST", url, body)
	req.Header.Set("Content-Type", contentType)

	return c.Send(req)
}

// DownloadFile performs simple file downloading and storing in a given path.
func (c *Client) DownloadFile(url, filename string) (int, error) {
	resp, err := c.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, errors.New(resp.Status)
	}

	file, closer, err := tools.CreateOrWriteFile(filename)
	if err != nil {
		return 0, err
	}

	if closer != nil {
		defer closer()
	}

	n, err := tools.WriteToFile(file, resp.Body)
	if err != nil {
		return 0, err
	}

	return n, nil
}
