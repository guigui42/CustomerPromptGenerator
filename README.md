# CustomerPromptGenerator

A Go CLI tool that generates customized demo prompts for creating landing pages on a Cat Supply Store, tailored to specific customer companies.

Built for GitHub Solutions Engineers to quickly create personalized demo content using the GitHub Copilot SDK.

## Prerequisites

- **GitHub Copilot CLI** installed and authenticated ([Installation guide](https://docs.github.com/en/copilot/how-tos/set-up/install-copilot-cli))
- **Go 1.21+** (for building from source)

## Installation

```bash
# Clone and build
cd CustomerPromptGenerator
go build -o CustomerPromptGenerator .
```

## Usage

```bash
./CustomerPromptGenerator "<company-name>"
```

### Example

```bash
./CustomerPromptGenerator "Spotify"
```

Output:
```
🔍 Researching Spotify...
✅ Found: https://www.spotify.com - Spotify is a digital music streaming...

📋 Generated Prompt (copied to clipboard):
──────────────────────────────────────────────────
Create a fun and engaging landing page for our Cat Supply Store...
[tailored prompt for Spotify]
──────────────────────────────────────────────────
✅ Prompt copied to clipboard!
```

The generated prompt is automatically copied to your clipboard, ready to paste into Copilot during your demo.

## Customizing the Prompt Template

Edit `prompt_template.txt` to customize the generated prompt. Available placeholders:

| Placeholder | Description |
|-------------|-------------|
| `{{.CompanyName}}` | The customer's company name |
| `{{.CompanyURL}}` | Their official website URL |
| `{{.CompanyDescription}}` | What the company does |
| `{{.Industry}}` | Their primary industry/sector |

If `prompt_template.txt` is not found, a sensible default template is used.

## How It Works

1. Takes a company name as input
2. Uses the GitHub Copilot SDK to research the company (URL, description, industry)
3. Fills in a customizable prompt template with the company details
4. Outputs the prompt to terminal and copies it to clipboard

## Demo Flow

1. Run the tool with your customer's company name
2. The prompt is automatically in your clipboard
3. Open your cat supply store demo repo (e.g., `octodemo/octocat_supply-crispy-pancake`)
4. Start Copilot and paste the generated prompt
5. Watch Copilot create a fun, personalized landing page!

## License

MIT
