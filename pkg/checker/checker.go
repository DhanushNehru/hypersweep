package checker

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/DhanushNehru/hypersweep/pkg/extractor"
)

// CheckResult adds the status to the extracted result
type CheckResult struct {
	Original extractor.Result
	Status   int
	IsAlive  bool
	Error    error
}

// Checker manages the HTTP worker pool
type Checker struct {
	client  *http.Client
	workers int
}

// NewChecker initializes a new checker with a configured realistic client
func NewChecker(workers int, timeout time.Duration) *Checker {
	// Create a transport that handles connection reuse and skips bad certs if needed (optional)
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return &Checker{
		client:  client,
		workers: workers,
	}
}

// Run checks all URLs concurrently and returns the results
func (c *Checker) Run(urls []extractor.Result) []CheckResult {
	jobs := make(chan extractor.Result, len(urls))
	results := make(chan CheckResult, len(urls))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go c.worker(jobs, results, &wg)
	}

	// Send jobs
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	close(results)

	// Collect results
	var finalResults []CheckResult
	for res := range results {
		finalResults = append(finalResults, res)
	}

	return finalResults
}

func (c *Checker) worker(jobs <-chan extractor.Result, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		results <- c.check(job)
	}
}

func (c *Checker) check(target extractor.Result) CheckResult {
	req, err := http.NewRequest("HEAD", target.URL, nil)
	if err != nil {
		return CheckResult{Original: target, IsAlive: false, Error: err}
	}

	// Spoof User-Agent as many generic sites block default gohttp
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	resp, err := c.client.Do(req)
	
	// Fast path for errors
	if err != nil {
		return CheckResult{Original: target, IsAlive: false, Error: err}
	}
	defer resp.Body.Close()

	// If HEAD request is rejected with 405 (Method Not Allowed) or 403 (Forbidden), we retry with GET
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusForbidden {
		reqGET, _ := http.NewRequest("GET", target.URL, nil)
		reqGET.Header.Set("User-Agent", req.Header.Get("User-Agent"))
		respGET, errGET := c.client.Do(reqGET)
		
		if errGET == nil {
			defer respGET.Body.Close()
			return CheckResult{
				Original: target,
				Status:   respGET.StatusCode,
				IsAlive:  respGET.StatusCode >= 200 && respGET.StatusCode < 400,
			}
		}
	}

	return CheckResult{
		Original: target,
		Status:   resp.StatusCode,
		IsAlive:  resp.StatusCode >= 200 && resp.StatusCode < 400,
	}
}
