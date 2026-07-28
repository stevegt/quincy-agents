// Intent: Parse markdown modules with selectable heading subtrees and stable
// agent_module metadata so presets can target blocks without depending on
// mutable heading text. Source: DI-lorad

package module

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var tagPattern = regexp.MustCompile(`\{([^}]+)\}`)

type Metadata struct {
	ID    string
	TLDR  string
	Group string
}

type Block struct {
	Type     string
	Level    int
	Heading  string
	Content  string
	Tags     []string
	Metadata Metadata
}

type Module struct {
	Name    string
	Path    string
	Content string
	Blocks  []Block
}

func Parse(path string) (*Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var blocks []Block
	content := string(data)
	lines := strings.Split(content, "\n")
	inCodeBlock := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		level, headerContent, ok := parseHeading(line)
		if !ok {
			continue
		}

		end := findHeadingSubtreeEnd(lines, i+1, level)
		blockLines := append([]string(nil), lines[i:end]...)
		metadata, err := extractHeadingMetadata(blockLines)
		if err != nil {
			return nil, fmt.Errorf("%s heading %q: %w", path, headerContent, err)
		}
		blocks = append(blocks, Block{
			Type:     "header",
			Level:    level,
			Heading:  strings.TrimSpace(tagPattern.ReplaceAllString(headerContent, "")),
			Content:  StripBuilderMetadata(strings.Join(blockLines, "\n")),
			Tags:     extractTags(headerContent),
			Metadata: metadata,
		})
	}

	return &Module{Path: path, Content: content, Blocks: blocks}, nil
}

func extractTags(header string) []string {
	matches := tagPattern.FindStringSubmatch(header)
	if len(matches) < 2 {
		return nil
	}

	tagStr := matches[1]
	var tags []string
	for _, t := range strings.Fields(tagStr) {
		t = strings.TrimPrefix(t, "#")
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func parseHeading(line string) (int, string, bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", false
	}

	level := 0
	for _, ch := range line {
		if ch == '#' {
			level++
		} else {
			break
		}
	}

	if level == 0 || level > 6 || len(line) <= level || line[level] != ' ' {
		return 0, "", false
	}

	return level, strings.TrimSpace(line[level:]), true
}

func findHeadingSubtreeEnd(lines []string, start int, parentLevel int) int {
	inCodeBlock := false
	for i := start; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}
		level, _, ok := parseHeading(line)
		if ok && level <= parentLevel {
			return i
		}
	}
	return len(lines)
}

func extractHeadingMetadata(blockLines []string) (Metadata, error) {
	var metadata Metadata
	if len(blockLines) < 2 {
		return metadata, nil
	}

	i := 1
	for i < len(blockLines) && strings.TrimSpace(blockLines[i]) == "" {
		i++
	}
	if i >= len(blockLines) || strings.TrimSpace(blockLines[i]) != "<!--" {
		return metadata, nil
	}

	var commentLines []string
	for ; i < len(blockLines); i++ {
		trimmed := strings.TrimSpace(blockLines[i])
		commentLines = append(commentLines, blockLines[i])
		if trimmed == "-->" {
			break
		}
	}
	if len(commentLines) == 0 || strings.TrimSpace(commentLines[len(commentLines)-1]) != "-->" {
		return metadata, fmt.Errorf("unterminated agent_module metadata comment")
	}

	return parseAgentModuleMetadata(commentLines), nil
}

func parseAgentModuleMetadata(commentLines []string) Metadata {
	var metadata Metadata
	inAgentModule := false
	for _, line := range commentLines {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "agent_module:":
			inAgentModule = true
		case trimmed == "<!--" || trimmed == "-->" || trimmed == "":
			continue
		case inAgentModule:
			key, value, ok := strings.Cut(trimmed, ":")
			if !ok {
				continue
			}
			value = strings.TrimSpace(value)
			value = strings.Trim(value, `"'`)
			switch strings.TrimSpace(key) {
			case "id":
				metadata.ID = value
			case "tldr":
				metadata.TLDR = value
			case "group":
				metadata.Group = value
			}
		}
	}
	return metadata
}

func StripBuilderMetadata(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inAgentModuleComment := false
	pendingComment := false
	var commentBuffer []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if pendingComment {
			commentBuffer = append(commentBuffer, line)
			if strings.Contains(trimmed, "agent_module:") {
				inAgentModuleComment = true
			}
			if trimmed == "-->" {
				if !inAgentModuleComment {
					result = append(result, commentBuffer...)
				}
				pendingComment = false
				inAgentModuleComment = false
				commentBuffer = nil
			}
			continue
		}

		if trimmed == "<!--" {
			pendingComment = true
			inAgentModuleComment = false
			commentBuffer = []string{line}
			continue
		}

		if strings.HasPrefix(line, "#") {
			line = stripTags(line)
		}
		result = append(result, line)
	}

	if pendingComment {
		result = append(result, commentBuffer...)
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}

func stripTags(line string) string {
	start := strings.Index(line, "{")
	end := strings.Index(line, "}")
	if start == -1 || end == -1 || end <= start {
		return line
	}
	return strings.TrimSpace(line[:start] + line[end+1:])
}
