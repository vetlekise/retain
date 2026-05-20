# Retain

A TUI for active learning and knowledge retention.

## Installation

**Using Go:**

```sh
go install github.com/vetlekise/retain@latest
```

**Pre-built binaries:**

Download the latest release for your platform from the [GitHub Releases](https://github.com/vetlekise/retain/releases) page.

## Usage

```bash
retain -exam <file.json>
```

### Keybinds

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `space` | Toggle option |
| `enter` | Submit answer / Next question |
| `→` | Next question |
| `p` / `←` | Previous question |
| `h` | Show hint |
| `q` / `ctrl+c` | Quit |

### Question bank format

Questions are defined in JSON. See [`examples/exam.json`](examples/exam.json) for a full example.

```json
{
  "meta": { "title": "My Exam", "version": "1.0.0" },
  "questions": [
    {
      "id": "q1",
      "type": "single",
      "text": "What is 2 + 2?",
      "options": [
        { "id": "a", "text": "3" },
        { "id": "b", "text": "4" }
      ],
      "correct_answers": ["b"],
      "explanation": "2 + 2 = 4.",
      "hint": "Think even numbers."
    }
  ]
}
```