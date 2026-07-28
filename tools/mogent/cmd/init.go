// Intent: Initialize block-native mogent configs and starter modules so new
// repos can immediately build AGENTS.md from selectable heading subtrees.
// Source: DI-soviv

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/internal/module"
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

	for _, category := range []string{"identity", "instructions", "constraints", "format"} {
		writeStarterModule(".mogent/modules", category)
		writeModuleConfig(config, category, category+".md", []string{category})
	}

	config.WriteString("[activate]\nscopes = []\n")
}

func writeInteractive(config *strings.Builder) {
	reader := bufio.NewReader(os.Stdin)
	moduleDir := ".mogent/modules"

	fmt.Println("Mogent builds your AGENTS.md from selectable Markdown heading blocks.")
	fmt.Println("Let's set up your block modules.")
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
		writeCategorySection(config, reader, cat, moduleDir)
	}

	config.WriteString("[activate]\nscopes = []\n")
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

func writeCategorySection(config *strings.Builder, reader *bufio.Reader, category string, moduleDir string) {
	fmt.Printf("\n--- %s ---\n", displayName(category))

	// Intent: Configure one block-native module per category, replacing the old
	// file/tag scan with immediate selectable block IDs. Source: DI-soviv
	source := category + ".md"
	path := filepath.Join(moduleDir, source)
	writeStarterModule(moduleDir, category)

	parsedModule, err := module.Parse(path)
	if err != nil {
		fmt.Printf("Could not inspect %s: %v\n", source, err)
		writeModuleConfig(config, category, source, []string{category})
		return
	}

	var blockIDs []string
	fmt.Printf("Module: %s\n", source)
	for _, block := range parsedModule.Blocks {
		if block.Metadata.ID == "" {
			continue
		}
		blockIDs = append(blockIDs, block.Metadata.ID)
		fmt.Printf("  - %-24s %s\n", block.Metadata.ID, block.Heading)
	}

	defaultBlocks := []string{category}
	fmt.Printf("Blocks [%s]: ", strings.Join(defaultBlocks, ","))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		defaultBlocks = splitCSV(input)
	}
	if len(blockIDs) == 0 {
		defaultBlocks = []string{category}
	}

	writeModuleConfig(config, category, source, defaultBlocks)
}

func writeModuleConfig(config *strings.Builder, category string, source string, blocks []string) {
	config.WriteString(fmt.Sprintf("[category.%s]\n", category))
	config.WriteString(fmt.Sprintf("[[category.%s.module]]\n", category))
	config.WriteString(fmt.Sprintf("name = \"%s\"\n", strings.TrimSuffix(source, ".md")))
	config.WriteString(fmt.Sprintf("source = \"%s\"\n", source))
	config.WriteString("blocks = [")
	for i, block := range blocks {
		if i > 0 {
			config.WriteString(", ")
		}
		config.WriteString(fmt.Sprintf("\"%s\"", block))
	}
	config.WriteString("]\n\n")
}

func writeStarterModule(moduleDir string, category string) {
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		fmt.Printf("Could not create %s: %v\n", moduleDir, err)
		return
	}

	path := filepath.Join(moduleDir, category+".md")
	if _, err := os.Stat(path); err == nil {
		return
	}

	content := fmt.Sprintf(`# %s

<!--
agent_module:
  id: %s
  tldr: Starter %s instructions.
-->
Add %s instructions here.
`, displayName(category), category, category, category)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Could not create %s: %v\n", path, err)
	}
}

func splitCSV(input string) []string {
	var values []string
	for _, value := range strings.Split(input, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func displayName(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '-' || r == '_' || r == '/'
	})
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
