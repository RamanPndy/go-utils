package goutils

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// In-cluster defaults.
	serviceAccountTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	serviceAccountCAFile    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	kubernetesServiceHost   = "KUBERNETES_SERVICE_HOST"
	kubernetesServicePort   = "KUBERNETES_SERVICE_PORT"
)

// Client is a minimal Kubernetes REST client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// --- K8s API list response types ---

type k8sListResponse struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Items      []json.RawMessage `json:"items"`
}

type k8sObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp string            `json:"creationTimestamp"`
}

type k8sObject struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   k8sObjectMeta   `json:"metadata"`
	Spec       json.RawMessage `json:"spec"`
	Status     json.RawMessage `json:"status"`
}

type k8sEvent struct {
	APIVersion     string        `json:"apiVersion"`
	Kind           string        `json:"kind"`
	Metadata       k8sObjectMeta `json:"metadata"`
	Type           string        `json:"type"`
	Reason         string        `json:"reason"`
	Message        string        `json:"message"`
	InvolvedObject struct {
		Kind      string `json:"kind"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		UID       string `json:"uid"`
	} `json:"involvedObject"`
	Source struct {
		Component string `json:"component"`
	} `json:"source"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
	Count          int    `json:"count"`
}

type k8sPodStatus struct {
	Phase             string          `json:"phase"`
	ContainerStatuses json.RawMessage `json:"containerStatuses"`
}

type k8sPod struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Metadata   k8sObjectMeta `json:"metadata"`
	Spec       struct {
		NodeName string `json:"nodeName"`
	} `json:"spec"`
	Status k8sPodStatus `json:"status"`
}

type k8sJob struct {
	Metadata k8sObjectMeta `json:"metadata"`
	Status   struct {
		StartTime      string `json:"startTime"`
		CompletionTime string `json:"completionTime"`
		Succeeded      int    `json:"succeeded"`
		Failed         int    `json:"failed"`
	} `json:"status"`
}

// KubectlApply runs `kubectl apply -f -` with the given manifest as stdin.
func KubectlApply(manifest string) error {
	cmd := exec.CommandContext(context.Background(), "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// KubectlGet runs `kubectl get <resource> [name] -o json` and returns the output.
func KubectlGet(resource, name string) ([]byte, error) {
	args := []string{"get", resource}
	if name != "" {
		args = append(args, name)
	}
	args = append(args, "-o", "json")
	return exec.CommandContext(context.Background(), "kubectl", args...).Output()
}

// KubectlRaw runs an arbitrary kubectl command and returns stdout.
func KubectlRaw(args ...string) ([]byte, error) {
	return exec.CommandContext(context.Background(), "kubectl", args...).Output()
}

// InClusterConfig builds a Client using the projected ServiceAccount token.
func InClusterConfig() (*Client, error) {
	host := os.Getenv(kubernetesServiceHost)
	port := os.Getenv(kubernetesServicePort)
	if host == "" || port == "" {
		return nil, fmt.Errorf("not running inside a Kubernetes pod (KUBERNETES_SERVICE_HOST/PORT not set)")
	}

	tokenBytes, err := os.ReadFile(serviceAccountTokenFile)
	if err != nil {
		return nil, fmt.Errorf("reading service account token: %w", err)
	}

	caCert, err := os.ReadFile(serviceAccountCAFile)
	if err != nil {
		return nil, fmt.Errorf("reading cluster CA: %w", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	tlsCfg := &tls.Config{RootCAs: pool}
	transport := &http.Transport{TLSClientConfig: tlsCfg}

	return &Client{
		baseURL:    fmt.Sprintf("https://%s:%s", host, port),
		httpClient: &http.Client{Transport: transport, Timeout: 30 * time.Second},
		token:      strings.TrimSpace(string(tokenBytes)),
	}, nil
}

// NewClientWithToken builds a Client using a provided base URL and bearer token.
// Useful for testing and local development.
func NewClientWithToken(baseURL, token string, skipTLSVerify bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify}, //nolint:gosec
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Transport: transport, Timeout: 30 * time.Second},
		token:      token,
	}
}

// Get performs a GET request against the given API path and decodes the response
// into out.
func (c *Client) Get(ctx context.Context, path string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s: HTTP %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// Apply creates or updates a resource using server-side apply (PATCH with
// application/apply-patch+yaml content type).  body should be a JSON/YAML
// manifest.
func (c *Client) Apply(ctx context.Context, path string, body []byte) error {
	url := c.baseURL + path + "?fieldManager=kubesleuth&force=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/apply-patch+yaml")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PATCH %s: HTTP %d", path, resp.StatusCode)
	}
	return nil
}

// Create POST-creates a resource.
func (c *Client) Create(ctx context.Context, path string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST %s: HTTP %d", path, resp.StatusCode)
	}
	return nil
}

// Delete removes a resource.
func (c *Client) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("DELETE %s: HTTP %d", path, resp.StatusCode)
	}
	return nil
}

// UpdateStatus PATCHes the /status subresource of a resource.
func (c *Client) UpdateStatus(ctx context.Context, path string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+path+"/status", bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/merge-patch+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PATCH status %s: HTTP %d", path, resp.StatusCode)
	}
	return nil
}

// setHeaders adds the Authorization and Accept headers to every request.
func (c *Client) setHeaders(req *http.Request) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")
}

func KubectlApplyManifest(manifest string) error {
	cmd := exec.CommandContext(context.Background(), "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(manifest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
