package report

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Vern discovery report</title>
<style>
:root {
  --bg: #fbfaf7; --fg: #1d1b18; --muted: #6b665e; --line: #e4e0d8; --card: #ffffff;
  --accent: #7a4fd1; --ok: #2f7d4a; --bad: #b3261e; --warn: #9a6400;
}
@media (prefers-color-scheme: dark) {
  :root { --bg: #151413; --fg: #ece9e3; --muted: #a39d93; --line: #2e2b27; --card: #1d1c1a;
    --accent: #b49cff; --ok: #6cc28a; --bad: #ff8a80; --warn: #f0b354; }
}
* { box-sizing: border-box; }
body { margin: 0; background: var(--bg); color: var(--fg);
  font: 15px/1.55 -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
main { max-width: 980px; margin: 0 auto; padding: 32px 16px 64px; }
h1 { font-size: 26px; margin: 0 0 4px; }
h2 { font-size: 19px; margin: 40px 0 12px; padding-bottom: 6px; border-bottom: 1px solid var(--line); }
.idea { margin: 8px 0 16px; }
.idea h1, .idea h2, .idea h3 { font-size: 15px; margin: 12px 0 4px; border: 0; padding: 0; }
.idea p { margin: 4px 0; }
.meta { color: var(--muted); font-size: 13px; display: flex; flex-wrap: wrap; gap: 6px 18px; }
.status { font-weight: 600; text-transform: uppercase; letter-spacing: .04em; }
.status.ok { color: var(--ok); } .status.failed { color: var(--bad); } .status.partial, .status.running { color: var(--warn); }
.table-wrap { overflow-x: auto; }
table { border-collapse: collapse; width: 100%; font-size: 13.5px; }
th, td { text-align: left; padding: 7px 10px; border-bottom: 1px solid var(--line); vertical-align: top; }
th { color: var(--muted); font-weight: 600; }
.tasks { display: grid; gap: 12px; }
.task { background: var(--card); border: 1px solid var(--line); border-radius: 10px; padding: 14px 16px; }
.task h3 { margin: 0 0 6px; font-size: 16px; }
.task .id { color: var(--accent); font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 13px; margin-right: 6px; }
.chip { display: inline-block; font-size: 12px; padding: 1px 8px; border-radius: 999px; border: 1px solid var(--line); color: var(--muted); margin-right: 6px; }
.task ul { margin: 6px 0 0 18px; padding: 0; }
.small { color: var(--muted); font-size: 13px; }
details { background: var(--card); border: 1px solid var(--line); border-radius: 10px; margin: 10px 0; }
summary { cursor: pointer; padding: 10px 14px; font-weight: 600; }
summary .small { font-weight: 400; margin-left: 8px; }
.doc { padding: 0 18px 14px; overflow-x: auto; }
.doc pre { background: var(--bg); border: 1px solid var(--line); border-radius: 6px; padding: 10px; overflow-x: auto; }
.doc code { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 13px; }
.doc table { margin: 8px 0; }
a { color: var(--accent); }
footer { margin-top: 48px; color: var(--muted); font-size: 12px; }
</style>
</head>
<body>
<main>
<h1>Vern discovery report</h1>
<div class="idea">{{.Idea}}</div>
{{with .Manifest}}
<div class="meta">
  {{if ne .Status "unknown"}}<span class="status {{.Status}}">{{.Status}}</span>{{end}}
  {{if .StartedAt}}<span>Started {{.StartedAt}}</span>{{end}}
  {{if duration .DurationMS}}<span>{{duration .DurationMS}}</span>{{end}}
  {{if .Pipeline}}<span>Pipeline: {{.Pipeline}}</span>{{end}}
  {{if .LLMMode}}<span>LLM mode: {{.LLMMode}}</span>{{end}}
  {{if .VernVersion}}<span>vern {{.VernVersion}}</span>{{end}}
  {{if .ReadOnly}}<span>read-only</span>{{end}}
  {{with .Usage}}<span>{{if .TotalTokens}}{{.TotalTokens}} tokens{{end}} {{cost .CostUSD}}</span>{{end}}
</div>

{{if .Steps}}
<h2>Pipeline</h2>
<div class="table-wrap"><table>
<tr><th>#</th><th>Step</th><th>Persona</th><th>LLM</th><th>Model</th><th>Status</th><th>Time</th><th>Usage</th></tr>
{{range .Steps}}
<tr>
  <td>{{.StepNum}}</td><td>{{.Name}}</td><td>{{.Persona}}</td>
  <td>{{if .FellBack}}<s>{{.OriginalLLM}}</s> → {{end}}{{.LLMUsed}}</td>
  <td>{{if .Model}}{{.Model}}{{else}}<span class="small">default</span>{{end}}{{if .Effort}} · {{.Effort}}{{end}}</td>
  <td><span class="status {{.Status}}">{{.Status}}</span>{{if .ErrorDetail}}<div class="small">{{.ErrorDetail}}</div>{{end}}</td>
  <td>{{duration .DurationMS}}</td>
  <td>{{with .Usage}}{{if .TotalTokens}}{{.TotalTokens}} tok {{end}}{{cost .CostUSD}}{{end}}</td>
</tr>
{{end}}
</table></div>
{{end}}
{{end}}

{{if .Tasks}}
<h2>Tasks ({{len .Tasks}})</h2>
<div class="tasks">
{{range .Tasks}}
<div class="task">
  <h3><span class="id">{{.ID}}</span>{{.Title}}</h3>
  <div>{{if .Complexity}}<span class="chip">{{.Complexity}}</span>{{end}}{{if .Dependencies}}<span class="chip">after {{join .Dependencies ", "}}</span>{{end}}</div>
  {{if .Description}}<p>{{.Description}}</p>{{end}}
  {{if .AcceptanceCriteria}}<ul>{{range .AcceptanceCriteria}}<li>{{.}}</li>{{end}}</ul>{{end}}
  {{if .Files}}<p class="small">Files: {{join .Files ", "}}</p>{{end}}
</div>
{{end}}
</div>
{{end}}

{{if .Oracle}}<h2>Oracle</h2>
<details open><summary>{{.Oracle.Title}}<span class="small">{{.Oracle.File}}</span></summary><div class="doc">{{.Oracle.HTML}}</div></details>
{{end}}

{{if or .Synth .Council}}<h2>VernHole</h2>
{{with .Synth}}<details open><summary>{{.Title}}<span class="small">{{.File}}</span></summary><div class="doc">{{.HTML}}</div></details>{{end}}
{{range .Council}}<details><summary>{{.Title}}<span class="small">{{.File}}</span></summary><div class="doc">{{.HTML}}</div></details>{{end}}
{{end}}

{{if .Steps}}<h2>Step outputs</h2>
{{range .Steps}}<details><summary>{{.Title}}<span class="small">{{.File}}</span></summary><div class="doc">{{.HTML}}</div></details>{{end}}
{{end}}

<footer>Built by <a href="https://github.com/jdonohoo/vern-bot">vern</a> on {{.Built}}.</footer>
</main>
</body>
</html>
`
