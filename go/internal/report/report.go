// Package report renders a discovery run as one self-contained HTML page.
package report

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/jdonohoo/vern-bot/go/internal/tasks"
)

// FileName is where Write saves the report, inside output/.
const FileName = "report.html"

// Section is one rendered markdown document.
type Section struct {
	Title string
	File  string
	HTML  template.HTML
}

type page struct {
	Manifest *pipeline.Manifest
	Idea     template.HTML
	Tasks    []tasks.Task
	Steps    []Section
	Council  []Section
	Synth    *Section
	Oracle   *Section
	Built    string
}

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

// render converts markdown to HTML. goldmark drops raw HTML by default, so
// model output cannot inject markup into the report.
func render(src []byte) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(string(src)))
	}
	return template.HTML(buf.String())
}

func section(title, path string) *Section {
	data, err := os.ReadFile(path)
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return &Section{Title: title, File: filepath.Base(path), HTML: render(data)}
}

// Build renders the report for a discovery folder.
func Build(discoveryDir string) ([]byte, error) {
	outputDir := filepath.Join(discoveryDir, "output")
	if _, err := os.Stat(outputDir); err != nil {
		return nil, fmt.Errorf("%s has no output/ folder: is it a discovery folder?", discoveryDir)
	}
	m, err := pipeline.ReadManifest(discoveryDir)
	if err != nil {
		// Runs from before run.json: show what the files say.
		m = &pipeline.Manifest{Kind: "discovery", Dir: discoveryDir, Status: "unknown"}
		if data, perr := os.ReadFile(filepath.Join(discoveryDir, "input", "prompt.md")); perr == nil {
			m.Idea = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(data)), "# Discovery Prompt"))
		}
	}
	p := page{Manifest: m, Idea: render([]byte(m.Idea)), Built: time.Now().Format("2006-01-02 15:04")}

	if doc, err := tasks.Load(outputDir); err == nil {
		p.Tasks = doc.Tasks
	}

	if len(m.Steps) > 0 {
		for _, s := range m.Steps {
			if sec := section(fmt.Sprintf("Step %d: %s", s.StepNum, s.Name), filepath.Join(discoveryDir, filepath.FromSlash(s.OutputFile))); sec != nil {
				p.Steps = append(p.Steps, *sec)
			}
		}
	} else {
		files, _ := filepath.Glob(filepath.Join(outputDir, "[0-9][0-9]-*.md"))
		sort.Strings(files)
		for _, f := range files {
			if sec := section(strings.TrimSuffix(filepath.Base(f), ".md"), f); sec != nil {
				p.Steps = append(p.Steps, *sec)
			}
		}
	}

	vh := filepath.Join(discoveryDir, "vernhole")
	p.Synth = section("VernHole synthesis", filepath.Join(vh, "synthesis.md"))
	council, _ := filepath.Glob(filepath.Join(vh, "[0-9][0-9]-*.md"))
	sort.Strings(council)
	for _, f := range council {
		if sec := section(strings.TrimSuffix(filepath.Base(f), ".md"), f); sec != nil {
			p.Council = append(p.Council, *sec)
		}
	}
	p.Oracle = section("Oracle vision", filepath.Join(discoveryDir, "oracle-vision.md"))
	if p.Oracle == nil {
		p.Oracle = section("Oracle vision", filepath.Join(vh, "oracle-vision.md"))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Write builds the report and saves it as output/report.html.
func Write(discoveryDir string) (string, error) {
	html, err := Build(discoveryDir)
	if err != nil {
		return "", err
	}
	path := filepath.Join(discoveryDir, "output", FileName)
	return path, os.WriteFile(path, html, 0644)
}

func duration(ms int64) string {
	if ms <= 0 {
		return ""
	}
	return (time.Duration(ms) * time.Millisecond).Round(time.Second).String()
}

var tmpl = template.Must(template.New("report").Funcs(template.FuncMap{
	"duration": duration,
	"join":     strings.Join,
	"cost": func(v float64) string {
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("$%.4f", v)
	},
}).Parse(pageTemplate))
