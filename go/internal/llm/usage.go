package llm

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// Usage is what a CLI reported about a run's tokens and cost. Fields a CLI
// does not report stay zero.
type Usage struct {
	InputTokens  int64   `json:"input_tokens,omitempty"`
	OutputTokens int64   `json:"output_tokens,omitempty"`
	TotalTokens  int64   `json:"total_tokens,omitempty"`
	CostUSD      float64 `json:"cost_usd,omitempty"`
}

// Add sums two usages.
func (u Usage) Add(o Usage) Usage {
	return Usage{
		InputTokens:  u.InputTokens + o.InputTokens,
		OutputTokens: u.OutputTokens + o.OutputTokens,
		TotalTokens:  u.TotalTokens + o.TotalTokens,
		CostUSD:      u.CostUSD + o.CostUSD,
	}
}

// IsZero reports whether nothing was reported.
func (u Usage) IsZero() bool { return u == Usage{} }

// claudeResult is the final object `claude -p --output-format json` prints.
type claudeResult struct {
	Type         string   `json:"type"`
	IsError      bool     `json:"is_error"`
	Result       *string  `json:"result"`
	TotalCostUSD *float64 `json:"total_cost_usd"`
	Usage        struct {
		InputTokens              int64 `json:"input_tokens"`
		OutputTokens             int64 `json:"output_tokens"`
		CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
		CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	} `json:"usage"`
}

// parseClaudeJSON unwraps claude's JSON output. When stdout is not the
// expected object it is returned unchanged, so a CLI change degrades to
// raw text instead of an empty answer.
func parseClaudeJSON(stdout, _ string) (text string, u Usage, failed bool) {
	start := strings.Index(stdout, "{")
	for start >= 0 {
		var r claudeResult
		if err := json.Unmarshal([]byte(strings.TrimSpace(stdout[start:])), &r); err == nil && r.Type == "result" && r.Result != nil {
			in := r.Usage.InputTokens + r.Usage.CacheReadInputTokens + r.Usage.CacheCreationInputTokens
			u = Usage{InputTokens: in, OutputTokens: r.Usage.OutputTokens, TotalTokens: in + r.Usage.OutputTokens}
			if r.TotalCostUSD != nil {
				u.CostUSD = *r.TotalCostUSD
			}
			return *r.Result, u, r.IsError
		}
		next := strings.Index(stdout[start+1:], "\n{")
		if next < 0 {
			break
		}
		start += next + 2
	}
	return stdout, Usage{}, false
}

var codexTokensRe = regexp.MustCompile(`(?i)tokens used[:\s]*([\d,]+)`)

// parseCodexUsage reads the "tokens used" line codex exec prints on stderr.
func parseCodexUsage(stdout, stderr string) (string, Usage, bool) {
	m := codexTokensRe.FindAllStringSubmatch(stderr, -1)
	if len(m) == 0 {
		return stdout, Usage{}, false
	}
	n, err := strconv.ParseInt(strings.ReplaceAll(m[len(m)-1][1], ",", ""), 10, 64)
	if err != nil {
		return stdout, Usage{}, false
	}
	return stdout, Usage{TotalTokens: n}, false
}
