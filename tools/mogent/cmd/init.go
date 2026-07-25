package cmd

import (
	"bufio"
	"fmt"
	"os"
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

	config.WriteString("[output]\n")
	config.WriteString("path = \"AGENTS.md\"\n")

	if err := os.WriteFile("AGENTS.toml", []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.toml: %w", err)
	}

	fmt.Println("Created AGENTS.toml")
	fmt.Println("Run 'mogent build' to assemble your AGENTS.md")

	return nil
}

func writeDefaults(config *strings.Builder) {
	config.WriteString("[order]\ncategories = [\"base\", \"coding\"]\n\n")

	config.WriteString("[category.base]\n")
	config.WriteString("[[category.base.module]]\n")
	config.WriteString("name = \"project-structure\"\n")
	config.WriteString("source = \"base/project-structure.md\"\n\n")

	config.WriteString("[category.coding]\n")
	config.WriteString("[[category.coding.module]]\n")
	config.WriteString("name = \"style\"\n")
	config.WriteString("source = \"coding/style.md\"\n\n")

	config.WriteString("[activate]\nscopes = []\n")
}

func writeInteractive(config *strings.Builder) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Mogent Init ===\n")

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
		writeCategorySection(config, reader, cat)
	}

	scopes := promptScopes(reader)
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
	fmt.Println("Categories define groups of modules (e.g., base, coding, team).")
	fmt.Println("Enter categories separated by commas, or press Enter for defaults:")
	fmt.Print("Categories [base,coding]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{"base", "coding"}
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

func writeCategorySection(config *strings.Builder, reader *bufio.Reader, category string) {
	fmt.Printf("\n--- Category: %s ---\n", category)

	config.WriteString(fmt.Sprintf("[category.%s]\n", category))

	tags := promptTags(reader, fmt.Sprintf("Tags for '%s' (comma-separated, or empty)", category))
	if len(tags) > 0 {
		config.WriteString("tags = [")
		for i, t := range tags {
			if i > 0 {
				config.WriteString(", ")
			}
			config.WriteString(fmt.Sprintf("\"%s\"", t))
		}
		config.WriteString("]\n")
	}

	config.WriteString("\n")

	for {
		fmt.Printf("Add module to '%s'? [Y/n]: ", category)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer == "n" || answer == "no" {
			break
		}

		writeModuleSection(config, reader, category)
	}
}

func writeModuleSection(config *strings.Builder, reader *bufio.Reader, category string) {
	fmt.Print("Module name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}

	fmt.Printf("Source path (relative to module_dir or absolute) [%s.md]: ", name)
	source, _ := reader.ReadString('\n')
	source = strings.TrimSpace(source)
	if source == "" {
		source = name + ".md"
	}

	config.WriteString(fmt.Sprintf("[[category.%s.module]]\n", category))
	config.WriteString(fmt.Sprintf("name = \"%s\"\n", name))
	config.WriteString(fmt.Sprintf("source = \"%s\"\n", source))

	tags := promptTags(reader, fmt.Sprintf("Tags for module '%s' (comma-separated, or empty)", name))
	if len(tags) > 0 {
		config.WriteString("tags = [")
		for i, t := range tags {
			if i > 0 {
				config.WriteString(", ")
			}
			config.WriteString(fmt.Sprintf("\"%s\"", t))
		}
		config.WriteString("]\n")
	}

	config.WriteString("\n")
}

func promptTags(reader *bufio.Reader, prompt string) []string {
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
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

func promptScopes(reader *bufio.Reader) []string {
	fmt.Println("\n--- Active Scopes ---")
	fmt.Println("Scopes determine which tagged sections are included.")
	fmt.Println("Examples: org/acme, lang/go, person/yourname")
	fmt.Println("Enter scopes separated by commas, or press Enter for none:")
	fmt.Print("Scopes: ")

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
