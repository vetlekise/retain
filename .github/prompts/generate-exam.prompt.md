---
description: "Generate a retain exam JSON file from a name and domain breakdown. Use when: generate exam, create exam, new exam, create questions, generate questions."
name: "Generate Exam"
argument-hint: "Paste the exam name followed by each domain and its percentage on separate lines"
agent: "agent"
---

Before generating anything, ask the user for the following if not already provided:

1. **Exam name** — the full display name (e.g. "Kubernetes and Cloud Native Security Associate")
2. **Exam slug** — short lowercase hyphenated filename (e.g. `kcsa`); suggest one based on the name and confirm
3. **Domain list** — each domain on its own line with its percentage weight, optionally followed by sub-topics indented beneath it. Example:
   ```
   Overview of Cloud Native Security 14%
     The 4Cs of Cloud Native Security
     Cloud Provider and Infrastructure Security
   Kubernetes Cluster Component Security 22%
     API Server
     Kubelet
   ```

Once you have all three, generate the exam JSON.

---

## Generation rules

Follow the structure in [examples/exam.json](../../examples/exam.json) exactly.

**`meta`**
- `title`: the exam name provided
- `description`: one sentence summarising the exam
- `tags`: derive from domain names and the exam slug

**`domains`** (top-level array, after `meta`)
- One entry per domain: `{ "name": "<domain name>", "percentage": <number>, "subdomains": ["...", "..."] }`
- `subdomains` is an array of strings listing the topics under that domain
- Order matches the input order
- The TUI uses this to display the domain breakdown screen

**`questions`**
- Generate **20 questions** total.
- Distribute questions across domains proportionally to their percentages, rounding to the nearest whole number so the total stays at 20. Add the domain name as a `"domain"` field on each question.
- Within each domain, cover the listed sub-topics where possible.
- Mix question types: lean on `single` and `multiple`, include at least one `true_false` and one `drag_drop` overall.
- Every question must have: `id`, `type`, `text`, `correct_answers`, `explanation`, `domain`.
- `options` required for `single`, `multiple`, `true_false`. For `true_false` options must be exactly `[{"id":"true","text":"True"},{"id":"false","text":"False"}]`.
- `drag_drop` uses `items`, `targets`, and `correct_answers` as `[{"item":"...","target":"..."}]` — no `options` field.
- Include 1–2 `hints` per question.
- Include at least one `links` entry per question pointing to real documentation.

**General**
- Do not include a `config` or `state` block.
- IDs: questions `q1`, `q2`, … — options `a`, `b`, `c`, … — drag-drop items/targets `i1`/`t1`, etc.

Output only the raw JSON. No explanation, no markdown fences.
Save the file to `exams/<slug>.json`.
