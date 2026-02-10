package pipeline

import (
	"fmt"

	"github.com/jdonohoo/vern-bot/go/internal/council"
)

// EventHandler receives pipeline progress events.
// Implementations control how progress is displayed (console, TUI, etc).
type EventHandler interface {
	OnStepStart(stepNum int, name string, llm string)
	OnStepComplete(stepNum int, result StepResult)
	OnStepRetry(stepNum int, attempt int, llm string)
	OnStepFallback(stepNum int, from string, to string)
	OnVernStart(idx int, vern council.Vern, llm string)
	OnVernComplete(idx int, result VernHoleResult)
	OnSynthesisStart(llm string)
	OnSynthesisComplete()
	OnPipelineComplete(results []StepResult)
	OnPipelineBanner(idea string, dir string, mode string, stepCount int, batchMode bool, resumeFrom int, maxRetries int, timeout int)
}

// ConsoleHandler is the default EventHandler that prints to stdout (current behavior).
type ConsoleHandler struct{}

func (h *ConsoleHandler) OnStepStart(stepNum int, name string, llm string) {
	// Handled inline in pipeline.go for now
}

func (h *ConsoleHandler) OnStepComplete(stepNum int, result StepResult) {
	// Handled inline in pipeline.go for now
}

func (h *ConsoleHandler) OnStepRetry(stepNum int, attempt int, llm string) {
	// Handled inline in pipeline.go for now
}

func (h *ConsoleHandler) OnStepFallback(stepNum int, from string, to string) {
	// Handled inline in pipeline.go for now
}

func (h *ConsoleHandler) OnVernStart(idx int, vern council.Vern, llm string) {
	fmt.Printf(">>> Vern %d: %s (%s)\n", idx+1, vern.Desc, llm)
}

func (h *ConsoleHandler) OnVernComplete(idx int, result VernHoleResult) {
	if result.Succeeded {
		fmt.Printf("    OK: %s\n", result.Vern.ID)
	} else {
		fmt.Printf("    FAILED: %s (exit %d)\n", result.Vern.ID, result.ExitCode)
	}
}

func (h *ConsoleHandler) OnSynthesisStart(llm string) {
	fmt.Printf(">>> Synthesizing with %s...\n", llm)
}

func (h *ConsoleHandler) OnSynthesisComplete() {
	fmt.Println("    Synthesis complete")
}

func (h *ConsoleHandler) OnPipelineComplete(results []StepResult) {
	completed := 0
	for _, r := range results {
		if r.Status == "ok" {
			completed++
		}
	}
	fmt.Printf("\nPipeline: %d/%d steps completed\n", completed, len(results))
}

func (h *ConsoleHandler) OnPipelineBanner(idea string, dir string, mode string, stepCount int, batchMode bool, resumeFrom int, maxRetries int, timeout int) {
	fmt.Println("=== VERN DISCOVERY PIPELINE ===")
	fmt.Printf("Idea: %s\n", idea)
	fmt.Printf("Discovery folder: %s\n", dir)
	fmt.Println("Output: Vern Task Spec (VTS)")
	fmt.Printf("Pipeline mode: %s (%d steps)\n", mode, stepCount)
	if batchMode {
		fmt.Println("Mode: batch (non-interactive)")
	}
	if resumeFrom > 0 {
		fmt.Printf("Resuming from: step %d\n", resumeFrom)
	}
	fmt.Printf("Max retries: %d\n", maxRetries)
	fmt.Printf("Timeout: %ds per step\n", timeout)
	fmt.Println()
}

// ChannelHandler sends events through a channel for TUI integration.
type ChannelHandler struct {
	Events chan Event
}

// Event is a pipeline event sent to the TUI.
type Event struct {
	Type string
	Data interface{}
}

// Event type constants.
const (
	EventStepStart        = "step_start"
	EventStepComplete     = "step_complete"
	EventStepRetry        = "step_retry"
	EventStepFallback     = "step_fallback"
	EventVernStart        = "vern_start"
	EventVernComplete     = "vern_complete"
	EventSynthesisStart   = "synthesis_start"
	EventSynthesisComplete = "synthesis_complete"
	EventPipelineComplete = "pipeline_complete"
	EventPipelineBanner   = "pipeline_banner"
)

type StepStartData struct {
	StepNum int
	Name    string
	LLM     string
}

type StepRetryData struct {
	StepNum int
	Attempt int
	LLM     string
}

type StepFallbackData struct {
	StepNum int
	From    string
	To      string
}

type VernStartData struct {
	Index int
	Vern  council.Vern
	LLM   string
}

type BannerData struct {
	Idea       string
	Dir        string
	Mode       string
	StepCount  int
	BatchMode  bool
	ResumeFrom int
	MaxRetries int
	Timeout    int
}

func NewChannelHandler() *ChannelHandler {
	return &ChannelHandler{
		Events: make(chan Event, 100),
	}
}

func (h *ChannelHandler) OnStepStart(stepNum int, name string, llm string) {
	h.Events <- Event{Type: EventStepStart, Data: StepStartData{StepNum: stepNum, Name: name, LLM: llm}}
}

func (h *ChannelHandler) OnStepComplete(stepNum int, result StepResult) {
	h.Events <- Event{Type: EventStepComplete, Data: result}
}

func (h *ChannelHandler) OnStepRetry(stepNum int, attempt int, llm string) {
	h.Events <- Event{Type: EventStepRetry, Data: StepRetryData{StepNum: stepNum, Attempt: attempt, LLM: llm}}
}

func (h *ChannelHandler) OnStepFallback(stepNum int, from string, to string) {
	h.Events <- Event{Type: EventStepFallback, Data: StepFallbackData{StepNum: stepNum, From: from, To: to}}
}

func (h *ChannelHandler) OnVernStart(idx int, vern council.Vern, llm string) {
	h.Events <- Event{Type: EventVernStart, Data: VernStartData{Index: idx, Vern: vern, LLM: llm}}
}

func (h *ChannelHandler) OnVernComplete(idx int, result VernHoleResult) {
	h.Events <- Event{Type: EventVernComplete, Data: result}
}

func (h *ChannelHandler) OnSynthesisStart(llm string) {
	h.Events <- Event{Type: EventSynthesisStart, Data: llm}
}

func (h *ChannelHandler) OnSynthesisComplete() {
	h.Events <- Event{Type: EventSynthesisComplete}
}

func (h *ChannelHandler) OnPipelineComplete(results []StepResult) {
	h.Events <- Event{Type: EventPipelineComplete, Data: results}
}

func (h *ChannelHandler) OnPipelineBanner(idea string, dir string, mode string, stepCount int, batchMode bool, resumeFrom int, maxRetries int, timeout int) {
	h.Events <- Event{Type: EventPipelineBanner, Data: BannerData{
		Idea: idea, Dir: dir, Mode: mode, StepCount: stepCount,
		BatchMode: batchMode, ResumeFrom: resumeFrom, MaxRetries: maxRetries, Timeout: timeout,
	}}
}

// Ensure both handlers satisfy the interface at compile time.
var _ EventHandler = (*ConsoleHandler)(nil)
var _ EventHandler = (*ChannelHandler)(nil)
