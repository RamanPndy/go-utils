package goutils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

type ContentType string

const (
	ApplicationJSON           ContentType = "application/json"
	ApplicationXML            ContentType = "application/xml"
	TextPlain                 ContentType = "text/plain"
	TextHTML                  ContentType = "text/html"
	ApplicationFormURLEncoded ContentType = "application/x-www-form-urlencoded"
)

type APIRequest struct {
	Method      HTTPMethod
	URL         string
	Headers     map[string]string
	Body        []byte
	QueryParams map[string]string
}

func NewAPIRequest() *APIRequest {
	return &APIRequest{}
}

func (r *APIRequest) SetMethod(method HTTPMethod) *APIRequest {
	r.Method = method
	return r
}

func (r *APIRequest) SetURL(url string) *APIRequest {
	r.URL = url
	return r
}

func (r *APIRequest) SetHeaders(headers map[string]string) *APIRequest {
	r.Headers = headers
	return r
}

func (r *APIRequest) SetBody(body []byte) *APIRequest {
	r.Body = body
	return r
}

func (r *APIRequest) SetQueryParams(queryParams map[string]string) *APIRequest {
	r.QueryParams = queryParams
	return r
}

func (r *APIRequest) AddHeader(key, value string) *APIRequest {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
	return r
}

func (r *APIRequest) AddQueryParam(key, value string) *APIRequest {
	if r.QueryParams == nil {
		r.QueryParams = make(map[string]string)
	}
	r.QueryParams[key] = value
	return r
}

func (r *APIRequest) SetContentType(contentType ContentType) *APIRequest {
	r.AddHeader("Content-Type", string(contentType))
	return r
}

func (r *APIRequest) SetAuthorization(token string) *APIRequest {
	r.AddHeader("Authorization", token)
	return r
}

func (r *APIRequest) SetBearerToken(token string) *APIRequest {
	r.AddHeader("Authorization", "Bearer "+token)
	return r
}

func (r *APIRequest) SetBasicAuth(username, password string) *APIRequest {
	r.AddHeader("Authorization", "Basic "+basicAuth(username, password))
	return r
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64Encode(auth)
}

func base64Encode(s string) string {
	return strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(s)), "=")
}

func (r *APIRequest) SetJSONBody(jsonBody []byte) *APIRequest {
	r.SetContentType(ApplicationJSON)
	r.Body = jsonBody
	return r
}

func (r *APIRequest) SetXMLBody(xmlBody []byte) *APIRequest {
	r.SetContentType(ApplicationXML)
	r.Body = xmlBody
	return r
}

func (r *APIRequest) SetFormURLEncodedBody(formData map[string]string) *APIRequest {
	r.SetContentType(ApplicationFormURLEncoded)
	formValues := make([]string, 0, len(formData))
	for key, value := range formData {
		formValues = append(formValues, key+"="+value)
	}
	r.Body = []byte(strings.Join(formValues, "&"))
	return r
}

func (r *APIRequest) SetPlainTextBody(text string) *APIRequest {
	r.SetContentType(TextPlain)
	r.Body = []byte(text)
	return r
}

func (r *APIRequest) SetHTMLBody(html string) *APIRequest {
	r.SetContentType(TextHTML)
	r.Body = []byte(html)
	return r
}

func (r *APIRequest) SetTimeout(timeoutSeconds int) *APIRequest {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers["Timeout"] = fmt.Sprintf("%d", timeoutSeconds)
	return r
}

func (r *APIRequest) GetMethod() HTTPMethod {
	return r.Method
}

func (r *APIRequest) GetURL() string {
	return r.URL
}

func (r *APIRequest) GetHeaders() map[string]string {
	return r.Headers
}

func (r *APIRequest) GetBody() []byte {
	return r.Body
}

func (r *APIRequest) GetQueryParams() map[string]string {
	return r.QueryParams
}

func (r *APIRequest) GetContentType() string {
	return GetContentTypeFromHeaders(r.Headers)
}

func (r *APIRequest) GetFullURL() string {
	if len(r.QueryParams) == 0 {
		return r.URL
	}
	return EncodeQueryParams(r.URL, r.QueryParams)
}

func (r *APIRequest) IsValid() bool {
	if !IsValidHTTPMethod(r.Method) {
		return false
	}
	if !IsValidURLWithQueryParams(r.URL, r.QueryParams) {
		return false
	}
	if r.Headers != nil {
		for key, value := range r.Headers {
			if !IsValidHeaderKey(key) || !IsValidHeaderValue(value) {
				return false
			}
		}
	}
	if !IsValidBody(r.Body) {
		return false
	}
	return true
}

func (r *APIRequest) String() string {
	var sb strings.Builder
	sb.WriteString(string(r.Method) + " " + r.GetFullURL() + "\n")
	for key, value := range r.Headers {
		sb.WriteString(key + ": " + value + "\n")
	}
	sb.WriteString("\n")
	sb.Write(r.Body)
	return sb.String()
}

type APIResult struct {
	StatusCode      int
	ContentType     string
	ContentLength   int64
	BodyBytes       []byte
	ResponseBody    io.ReadCloser
	ResponseHeaders http.Header
}

type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

type API interface {
	DoRequest(ctx context.Context, request *APIRequest) (*APIResult, *APIError)
}

type APIClient struct {
}

func NewAPIClient() *APIClient {
	return &APIClient{}
}

func (c *APIClient) DoRequest(ctx context.Context, request *APIRequest) (*APIResult, *APIError) {
	switch request.Method {
	case GET:
		return c.do(ctx, http.MethodGet, request)
	case POST:
		return c.do(ctx, http.MethodPost, request)
	case PUT:
		return c.do(ctx, http.MethodPut, request)
	case PATCH:
		return c.do(ctx, http.MethodPatch, request)
	case DELETE:
		return c.do(ctx, http.MethodDelete, request)
	case HEAD:
		return c.do(ctx, http.MethodHead, request)
	case OPTIONS:
		return c.do(ctx, http.MethodOptions, request)
	default:
		return nil, &APIError{Message: "Unsupported HTTP method"}
	}
}

func (c *APIClient) do(ctx context.Context, requestMethod string, request *APIRequest) (*APIResult, *APIError) {
	req, reqErr := http.NewRequestWithContext(ctx, requestMethod, request.GetFullURL(), nil)
	if reqErr != nil {
		return nil, &APIError{Message: fmt.Sprintf("build request: %v", reqErr)}
	}
	resp, doErr := http.DefaultClient.Do(req)
	if doErr != nil {
		return nil, &APIError{Message: fmt.Sprintf("do request: %v", doErr)}
	}

	defer resp.Body.Close()
	result, err := c.constructResult(request, resp)
	if err != nil {
		return nil, &APIError{Message: fmt.Sprintf("construct result: %v", err)}
	}
	return result, nil
}

func (c *APIClient) DoRequestWithRetries(ctx context.Context, request *APIRequest, retries int) (*APIResult, *APIError) {
	var lastErr *APIError
	for i := 0; i <= retries; i++ {
		result, err := c.DoRequest(ctx, request)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *APIClient) DoRequestWithTimeout(ctx context.Context, request *APIRequest, timeoutSeconds int) (*APIResult, *APIError) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	return c.DoRequest(timeoutCtx, request)
}

func (c *APIClient) DoRequestWithRetriesAndTimeout(ctx context.Context, request *APIRequest, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	var lastErr *APIError
	for i := 0; i <= retries; i++ {
		result, err := c.DoRequest(timeoutCtx, request)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *APIClient) DoRequestWithCustomClient(ctx context.Context, request *APIRequest, client *http.Client) (*APIResult, *APIError) {
	req, reqErr := http.NewRequestWithContext(ctx, string(request.Method), request.GetFullURL(), nil)
	if reqErr != nil {
		return nil, &APIError{Message: fmt.Sprintf("build request: %v", reqErr)}
	}
	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, &APIError{Message: fmt.Sprintf("do request: %v", doErr)}
	}

	defer resp.Body.Close()

	result, err := c.constructResult(request, resp)
	if err != nil {
		return nil, &APIError{Message: fmt.Sprintf("construct result: %v", err)}
	}
	return result, nil
}

func (c *APIClient) constructResult(request *APIRequest, response *http.Response) (*APIResult, error) {
	bodyBytes, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read response body: %v", readErr)
	}

	result := &APIResult{
		StatusCode:      response.StatusCode,
		ContentType:     GetContentTypeFromHeaders(request.Headers),
		ContentLength:   response.ContentLength,
		BodyBytes:       bodyBytes,
		ResponseBody:    response.Body,
		ResponseHeaders: response.Header,
	}
	return result, nil
}

func (c *APIClient) DoRequestWithCustomClientAndRetries(ctx context.Context, request *APIRequest, client *http.Client, retries int) (*APIResult, *APIError) {
	var lastErr *APIError
	for i := 0; i <= retries; i++ {
		result, err := c.DoRequestWithCustomClient(ctx, request, client)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *APIClient) DoRequestWithCustomClientAndTimeout(ctx context.Context, request *APIRequest, client *http.Client, timeoutSeconds int) (*APIResult, *APIError) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	return c.DoRequestWithCustomClient(timeoutCtx, request, client)
}

func (c *APIClient) DoRequestWithCustomClientRetriesAndTimeout(ctx context.Context, request *APIRequest, client *http.Client, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	var lastErr *APIError
	for i := 0; i <= retries; i++ {
		result, err := c.DoRequestWithCustomClient(timeoutCtx, request, client)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *APIClient) DoRequestWithCustomClientRetriesTimeoutAndHeaders(ctx context.Context, request *APIRequest, client *http.Client, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	var lastErr *APIError
	for i := 0; i <= retries; i++ {
		result, err := c.DoRequestWithCustomClient(timeoutCtx, request, client)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *APIClient) IsURLReachable(ctx context.Context, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func (c *APIClient) IsURLReachableWithTimeout(ctx context.Context, url string, timeoutSeconds int) bool {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func (c *APIClient) IsURLReachableWithCustomClient(ctx context.Context, url string, client *http.Client) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func (c *APIClient) EnableHTTPDebugging() {
	http.DefaultTransport = &debugRoundTripper{rt: http.DefaultTransport}
}

type debugRoundTripper struct {
	rt http.RoundTripper
}

func (d *debugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("Request: %s %s\n", req.Method, req.URL)
	for key, value := range req.Header {
		fmt.Printf("Header: %s: %s\n", key, strings.Join(value, ", "))
	}
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			fmt.Printf("Body: %s\n", string(bodyBytes))
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}
	}
	resp, err := d.rt.RoundTrip(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Response Status: %s\n", resp.Status)
	for key, value := range resp.Header {
		fmt.Printf("Response Header: %s: %s\n", key, strings.Join(value, ", "))
	}
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Printf("Response Body: %s\n", string(bodyBytes))
			resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}
	}
	return resp, nil
}

func (c *APIClient) DoRequestWithDebugging(ctx context.Context, request *APIRequest) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequest(ctx, request)
}

func (c *APIClient) DoRequestWithCustomClientAndDebugging(ctx context.Context, request *APIRequest, client *http.Client) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequestWithCustomClient(ctx, request, client)
}

func (c *APIClient) DoRequestWithRetriesTimeoutAndDebugging(ctx context.Context, request *APIRequest, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequestWithRetriesAndTimeout(ctx, request, retries, timeoutSeconds)
}

func (c *APIClient) DoRequestWithCustomClientRetriesTimeoutAndDebugging(ctx context.Context, request *APIRequest, client *http.Client, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequestWithCustomClientRetriesAndTimeout(ctx, request, client, retries, timeoutSeconds)
}

func (c *APIClient) DoRequestWithCustomClientRetriesTimeoutHeadersAndDebugging(ctx context.Context, request *APIRequest, client *http.Client, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequestWithCustomClientRetriesTimeoutAndHeaders(ctx, request, client, retries, timeoutSeconds)
}

func (c *APIClient) DoRequestWithCustomClientRetriesTimeoutHeadersAndDebuggingAndCustomErrorHandling(ctx context.Context, request *APIRequest, client *http.Client, retries int, timeoutSeconds int) (*APIResult, *APIError) {
	c.EnableHTTPDebugging()
	return c.DoRequestWithCustomClientRetriesTimeoutHeadersAndDebugging(ctx, request, client, retries, timeoutSeconds)
}

func (c *APIClient) SetCustomHTTPClient(client *http.Client) {
	http.DefaultClient = client
}

func (c *APIClient) SetCustomHTTPTransport(transport http.RoundTripper) {
	http.DefaultTransport = transport
}

func (c *APIClient) SetCustomHTTPTimeout(timeoutSeconds int) {
	http.DefaultClient.Timeout = time.Duration(timeoutSeconds) * time.Second
}

func (c *APIClient) SetCustomHTTPHeaders(headers map[string]string) {
	http.DefaultClient.Transport = &headerRoundTripper{rt: http.DefaultTransport, headers: headers}
}

type headerRoundTripper struct {
	rt      http.RoundTripper
	headers map[string]string
}

func (h *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.headers {
		req.Header.Set(key, value)
	}
	return h.rt.RoundTrip(req)
}

func (c *APIClient) SetCustomHTTPClientHeaders(client *http.Client, headers map[string]string) {
	client.Transport = &headerRoundTripper{rt: client.Transport, headers: headers}
}

func (c *APIClient) SetCustomHTTPClientTimeout(client *http.Client, timeoutSeconds int) {
	client.Timeout = time.Duration(timeoutSeconds) * time.Second
}

func (c *APIClient) SetCustomHTTPClientTransport(client *http.Client, transport http.RoundTripper) {
	client.Transport = transport
}

func (c *APIClient) SetCORSHeaders(headers map[string]string) {
	http.DefaultClient.Transport = &corsRoundTripper{rt: http.DefaultTransport, headers: headers}
}

type corsRoundTripper struct {
	rt      http.RoundTripper
	headers map[string]string
}

func (c *corsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	return c.rt.RoundTrip(req)
}

func EncodeQueryParams(url string, keyValue map[string]string) string {
	queryParams := make([]string, 0, len(keyValue))
	for key, value := range keyValue {
		queryParams = append(queryParams, key+"="+value)
	}
	if !strings.HasSuffix(url, "?") {
		url += "?"
	}
	return url + strings.Join(queryParams, "&")
}

func DecodeQueryParams(url string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(url, "?")
	if len(parts) < 2 {
		return result
	}
	queryString := parts[1]
	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

func IsValidURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

func IsValidHTTPMethod(method HTTPMethod) bool {
	switch method {
	case GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS:
		return true
	default:
		return false
	}
}

func IsValidContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	switch contentType {
	case string(ApplicationJSON), string(ApplicationXML), string(TextPlain), string(TextHTML), string(ApplicationFormURLEncoded):
		return true
	default:
		return false
	}
}

func IsValidStatusCode(statusCode int) bool {
	return statusCode >= 100 && statusCode <= 599
}

func IsValidHeaderKey(key string) bool {
	if key == "" {
		return false
	}
	for _, r := range key {
		if r < 33 || r > 126 || r == ':' {
			return false
		}
	}
	return true
}

func IsValidHeaderValue(value string) bool {
	for _, r := range value {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func IsValidBody(body []byte) bool {
	return body != nil
}

func IsValidQueryParamKey(key string) bool {
	if key == "" {
		return false
	}
	for _, r := range key {
		if r < 33 || r > 126 || r == '=' || r == '&' {
			return false
		}
	}
	return true
}

func IsValidQueryParamValue(value string) bool {
	for _, r := range value {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func IsValidURLWithQueryParams(url string, queryParams map[string]string) bool {
	if !IsValidURL(url) {
		return false
	}
	for key, value := range queryParams {
		if !IsValidQueryParamKey(key) || !IsValidQueryParamValue(value) {
			return false
		}
	}
	return true
}

func IsValidAPIResult(result APIResult) bool {
	if !IsValidStatusCode(result.StatusCode) {
		return false
	}
	if !IsValidContentType(result.ContentType) {
		return false
	}
	if !IsValidBody(result.BodyBytes) {
		return false
	}
	return true
}

func GetContentTypeFromHeaders(headers map[string]string) string {
	for key, value := range headers {
		if strings.EqualFold(key, "Content-Type") {
			return value
		}
	}
	return ""
}
