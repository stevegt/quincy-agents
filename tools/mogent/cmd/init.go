package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mogent in the current directory",
	Long:  "Create an AGENTS.toml configuration file with interactive questionnaire",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().Bool("defaults", false, "Use default configuration without prompts")
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat("AGENTS.toml"); err == nil {
		fmt.Println("AGENTS.toml already exists. Overwrite? [y/N]")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	useDefaults, _ := cmd.Flags().GetBool("defaults")

	var config strings.Builder

	config.WriteString("[config]\n")
	config.WriteString("module_dir = \".mogent/modules\"\n\n")

	if useDefaults {
		writeDefaults(&config)
	} else {
		writeInteractive(&config)
	}

	config.WriteString("\n[output]\n")
	config.WriteString("path = \"AGENTS.md\"\n")

	if err := os.WriteFile("AGENTS.toml", []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.toml: %w", err)
	}

	fmt.Println("\nCreated AGENTS.toml")
	fmt.Println("Run 'mogent build' to assemble your AGENTS.md")

	return nil
}

func writeDefaults(config *strings.Builder) {
	config.WriteString("[order]\ncategories = [\"identity\", \"instructions\", \"constraints\", \"format\"]\n\n")

	config.WriteString("[category.identity]\n")
	config.WriteString("[[category.identity.module]]\n")
	config.WriteString("name = \"identity\"\n")
	config.WriteString("source = \"identity.md\"\n\n")

	config.WriteString("[category.instructions]\n")
	config.WriteString("[[category.instructions.module]]\n")
	config.WriteString("name = \"instructions\"\n")
	config.WriteString("source = \"instructions.md\"\n\n")

	config.WriteString("[category.constraints]\n")
	config.WriteString("[[category.constraints.module]]\n")
	config.WriteString("name = \"constraints\"\n")
	config.WriteString("source = \"constraints.md\"\n\n")

	config.WriteString("[category.format]\n")
	config.WriteString("[[category.format.module]]\n")
	config.WriteString("name = \"format\"\n")
	config.WriteString("source = \"format.md\"\n\n")

	config.WriteString("[activate]\nscopes = []\n")
}

func writeInteractive(config *strings.Builder) {
	reader := bufio.NewReader(os.Stdin)
	moduleDir := ".mogent/modules"
	knownTags := make(map[string]bool)

	fmt.Println("Mogent builds your AGENTS.md from smaller Markdown files.")
	fmt.Println("Let's set up your sections and files.")
	fmt.Println()

	categories := promptCategoryList(reader)

	config.WriteString("[order]\ncategories = [")
	for i, cat := range categories {
		if i > 0 {
			config.WriteString(", ")
		}
		config.WriteString(fmt.Sprintf("\"%s\"", cat))
	}
	config.WriteString("]\n\n")

	for _, cat := range categories {
		writeCategorySection(config, reader, cat, moduleDir, knownTags)
	}

	scopes := promptScopes(reader, knownTags)
	config.WriteString("[activate]\nscopes = [")
	for i, scope := range scopes {
		if i > 0 {
			config.WriteString(", ")
		}
		config.WriteString(fmt.Sprintf("\"%s\"", scope))
	}
	config.WriteString("]\n")
}

func promptCategoryList(reader *bufio.Reader) []string {
	fmt.Println("What sections should your AGENTS.md have?")
	fmt.Println("Common choices: identity, instructions, constraints, format")
	fmt.Println("Press Enter for defaults, or type your own separated by commas.")
	fmt.Print("\nSections [identity,instructions,constraints,format]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{"identity", "instructions", "constraints", "format"}
	}

	var cats []string
	for _, c := range strings.Split(input, ",") {
		c = strings.TrimSpace(c)
		if c != "" {
			cats = append(cats, c)
		}
	}
	return cats
}

func writeCategorySection(config *strings.Builder, reader *bufio.Reader, category string, moduleDir string, knownTags map[string]bool) {
	fmt.Printf("\n--- %s ---\n", strings.Title(category))

	config.WriteString(fmt.Sprintf("[category.%s]\n", category))

	// Scan for existing files in this category
	existingFiles := findModuleFiles(filepath.Join(moduleDir, category))
	if len(existingFiles) > 0 {
		fmt.Println("Existing files:")
		for _, f := range existingFiles {
			fmt.Printf("  - %s\n", filepath.Base(f))
		}
	} else {
		fmt.Println("No files found yet.")
	}

	tags := promptTags(reader, "Only include this section when these tags are active (comma-separated, or empty for always)", knownTags)
	if len(tags) > 0 {
		config.WriteString("tags = [")
		for i, t := range tags {
			if i > 0 {
				config.WriteString(", ")
			}
			config.WriteString(fmt.Sprintf("\"%s\"", t))
			knownTags[t] = true
		}
		config.WriteString("]\n")
	}

	config.WriteString("\n")

	for {
		fmt.Printf("Add a file to '%s'? [Y/n]: ", category)
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "n" || input == "no" {
			break
		}
		if input != "" && input != "y" && input != "yes" {
			fmt.Println("  Please enter y or n.")
			continue
		}

		fmt.Print("File name (e.g., 'coding-style'): ")
		name, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Println("  (skipped)")
			continue
		}

		defaultSource := name + ".md"
		fmt.Printf("Path [%s]: ", defaultSource)
		source, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		source = strings.TrimSpace(source)
		if source == "" {
			source = defaultSource
		}

		config.WriteString(fmt.Sprintf("[[category.%s.module]]\n", category))
		config.WriteString(fmt.Sprintf("name = \"%s\"\n", name))
		config.WriteString(fmt.Sprintf("source = \"%s\"\n", source))

		// Show tags from the file if it exists
		filePath := filepath.Join(moduleDir, category, source)
		fileTags := extractTagsFromFile(filePath)
		if len(fileTags) > 0 {
			fmt.Printf("  Tags found in file: %s\n", strings.Join(fileTags, ", "))
		}

		tags := promptTags(reader, fmt.Sprintf("Only include '%s' when these tags are active (comma-separated, or empty)", name), knownTags)
		if len(tags) > 0 {
			config.WriteString("tags = [")
			for i, t := range tags {
				if i > 0 {
					config.WriteString(", ")
				}
				config.WriteString(fmt.Sprintf("\"%s\"", t))
				knownTags[t] = true
			}
			config.WriteString("]\n")
		}

		config.WriteString("\n")
	}
}

func promptTags(reader *bufio.Reader, prompt string, knownTags map[string]bool) []string {
	if len(knownTags) > 0 {
		fmt.Println("  Existing tags:")
		printTagTree(knownTags)
	}

	fmt.Printf("%s: ", prompt)
	input, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

	var tags []string
	for _, t := range strings.Split(input, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func promptScopes(reader *bufio.Reader, knownTags map[string]bool) []string {
	fmt.Println("\n--- Active Scopes ---")
	fmt.Println("Scopes filter which sections appear in your AGENTS.md.")
	fmt.Println("For example, 'lang/go' only includes sections tagged with 'lang/go'.")
	fmt.Println("Leave empty to include everything.")

	if len(knownTags) > 0 {
		fmt.Println("\nAvailable tags:")
		printTagTree(knownTags)
	}

	fmt.Print("\nActive scopes (comma-separated, or empty): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

	var scopes []string
	for _, s := range strings.Split(input, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			scopes = append(scopes, s)
		}
	}
	return scopes
}
