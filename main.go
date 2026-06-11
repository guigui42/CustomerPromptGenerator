package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	copilot "github.com/github/copilot-sdk/go"
)

// CompanyInfo holds the researched company data
type CompanyInfo struct {
	CompanyName        string
	CompanyURL         string
	CompanyDescription string
	Industry           string
}

// Default prompt template embedded in the binary
const defaultTemplate = `Create a fun and engaging landing page for our Cat Supply Store (OctoCat Supply) that cross-sells cat products to {{.CompanyName}} employees and customers.

## Context about {{.CompanyName}}:
- Website: {{.CompanyURL}}
- What they do: {{.CompanyDescription}}
- Industry: {{.Industry}}

## Your task:
Design a landing page that playfully connects {{.CompanyName}}'s business with cat ownership. The page should:

1. **Headline**: Create a clever, punny headline that relates {{.CompanyName}}'s industry to cats (e.g., if they're in music streaming, something like "Your cats deserve a purrfect playlist of products")

2. **Hero Section**: Write engaging copy that speaks to {{.CompanyName}} employees/customers as cat owners, making relatable jokes about their industry

3. **Featured Products**: Suggest 3-4 cat products that humorously tie into what {{.CompanyName}} does:
   - Each product should have a witty name related to their business
   - Include a fun description connecting the product to their industry
   - Keep it lighthearted and professional

4. **Call-to-Action**: A clever CTA that references their business

5. **Tone**: Professional but playful, showing we understand both their business AND cat owners

**Important**: Create a brand NEW, standalone landing page (e.g. a new file/route). Do NOT modify or overwrite the existing landing page or home page.

Make it memorable and shareable - something that would make a {{.CompanyName}} employee smile and want to show their coworkers!

Take a screenshot of the new landing page once created.

Create a draft GitHub issue with the full specifications.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: CustomerPromptGenerator <company-name>")
		fmt.Println("Example: CustomerPromptGenerator \"Spotify\"")
		os.Exit(1)
	}

	companyName := strings.Join(os.Args[1:], " ")
	fmt.Printf("🔍 Researching %s...\n", companyName)

	// Research the company using Copilot SDK
	info, err := researchCompany(companyName)
	if err != nil {
		fmt.Printf("❌ Error researching company: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Found: %s - %s\n\n", info.CompanyURL, truncate(info.CompanyDescription, 60))

	// Generate the prompt from template
	prompt, err := generatePrompt(info)
	if err != nil {
		fmt.Printf("❌ Error generating prompt: %v\n", err)
		os.Exit(1)
	}

	// Print the prompt
	fmt.Println("📋 Generated Prompt (copied to clipboard):")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(prompt)
	fmt.Println(strings.Repeat("─", 50))

	// Copy to clipboard
	if err := copyToClipboard(prompt); err != nil {
		fmt.Printf("⚠️  Could not copy to clipboard: %v\n", err)
	} else {
		fmt.Println("✅ Prompt copied to clipboard!")
	}
}

func researchCompany(companyName string) (*CompanyInfo, error) {
	ctx := context.Background()

	client := copilot.NewClient(nil)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start Copilot client: %w", err)
	}
	defer client.Stop()

	session, err := client.CreateSession(ctx, &copilot.SessionConfig{
		Model: "claude-haiku-4.5",
		OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Ask Copilot to research the company
	prompt := fmt.Sprintf(`Research the company "%s" and provide the following information in a structured format:

1. COMPANY_URL: The official website URL
2. DESCRIPTION: A 1-2 sentence description of what the company does
3. INDUSTRY: The primary industry/sector they operate in

Respond ONLY in this exact format (no markdown, no extra text):
COMPANY_URL: <url>
DESCRIPTION: <description>
INDUSTRY: <industry>`, companyName)

	response, err := session.SendPromptAndWait(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("empty response from Copilot")
	}

	// Extract content from the response
	var responseContent string
	if data, ok := response.Data.(*copilot.AssistantMessageData); ok {
		responseContent = data.Content
	}

	if responseContent == "" {
		return nil, fmt.Errorf("empty response content from Copilot")
	}

	// Parse the response
	info := parseCompanyResponse(responseContent, companyName)
	return info, nil
}

func parseCompanyResponse(response string, companyName string) *CompanyInfo {
	info := &CompanyInfo{
		CompanyName:        companyName,
		CompanyURL:         "https://www.google.com/search?q=" + strings.ReplaceAll(companyName, " ", "+"),
		CompanyDescription: "A company",
		Industry:           "Technology",
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "COMPANY_URL:") {
			info.CompanyURL = strings.TrimSpace(strings.TrimPrefix(line, "COMPANY_URL:"))
		} else if strings.HasPrefix(line, "DESCRIPTION:") {
			info.CompanyDescription = strings.TrimSpace(strings.TrimPrefix(line, "DESCRIPTION:"))
		} else if strings.HasPrefix(line, "INDUSTRY:") {
			info.Industry = strings.TrimSpace(strings.TrimPrefix(line, "INDUSTRY:"))
		}
	}

	return info
}

func generatePrompt(info *CompanyInfo) (string, error) {
	// Try to load custom template, fall back to default
	tmplContent := defaultTemplate
	if customTemplate, err := os.ReadFile("prompt_template.txt"); err == nil {
		tmplContent = string(customTemplate)
	}

	tmpl, err := template.New("prompt").Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, info); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard tool found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
