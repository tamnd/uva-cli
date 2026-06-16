// Package uva is the library behind the uva command line:
// the HTTP client, request shaping, and typed data models for UVa Online Judge
// problems fetched via the uHunt API.
package uva

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (compatible; uva-cli/0.1; +https://github.com/tamnd/uva-cli)"

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://uhunt.onlinejudge.org",
		UserAgent: DefaultUserAgent,
		Rate:      500 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the uHunt API for UVa Online Judge problems.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get: %w", lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// List fetches all problems from the uHunt API.
func (c *Client) List(ctx context.Context, limit int) ([]Problem, error) {
	raw, err := c.get(ctx, c.cfg.BaseURL+"/api/p")
	if err != nil {
		return nil, err
	}
	var items [][]any
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return wireToProblems(items, limit), nil
}

// Search searches problems by title (client-side filter).
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Problem, error) {
	all, err := c.List(ctx, 0)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Problem
	rank := 0
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Title), q) {
			rank++
			p.Rank = rank
			out = append(out, p)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

// GetByNum fetches a single problem by its problem number.
func (c *Client) GetByNum(ctx context.Context, num int) (*Problem, error) {
	url := fmt.Sprintf("%s/api/p/num/%d", c.cfg.BaseURL, num)
	raw, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	var item []any
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, err
	}
	if len(item) < 21 {
		return nil, nil
	}
	probs := wireToProblems([][]any{item}, 1)
	if len(probs) == 0 {
		return nil, nil
	}
	return &probs[0], nil
}

func wireToProblems(raw [][]any, limit int) []Problem {
	var out []Problem
	rank := 0
	for _, item := range raw {
		if len(item) < 21 {
			continue
		}
		status, _ := item[20].(float64)
		if status != 1 {
			continue
		}
		rank++
		if limit > 0 && rank > limit {
			break
		}
		pid := int(item[0].(float64))
		num := int(item[1].(float64))
		title, _ := item[2].(string)
		dacu := int(item[3].(float64))
		tle := int(item[14].(float64))
		wa := int(item[16].(float64))
		ac := int(item[18].(float64))
		rtl := int(item[19].(float64))
		out = append(out, Problem{
			Rank:        rank,
			PID:         pid,
			Num:         num,
			Title:       title,
			AC:          ac,
			WA:          wa,
			TLE:         tle,
			DACU:        dacu,
			TimeLimitMS: rtl,
			URL: fmt.Sprintf(
				"https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&page=show_problem&problem=%d",
				pid,
			),
		})
	}
	return out
}
