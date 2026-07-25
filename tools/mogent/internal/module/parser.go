// Intent: Parse markdown modules with header attribute tags {#tag1 #tag2}.
// Source: DI-jusuk

package module

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var tagPattern = regexp.MustCompile(`\{([^}]+)\}`)

type Block struct {
	Type    string
	Level   int
	Content string
	Tags    []string
}

type Module struct {
	Name   string
	Path   string
	Blocks []Block
}

func Parse(path string) (*Module, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []Block
	scanner := bufio.NewScanner(f)
	var currentBlock Block
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				currentBlock.Content += line + "\n"
				blocks = append(blocks, currentBlock)
				currentBlock = Block{}
				inCodeBlock = false
				continue
			} else {
				inCodeBlock = true
				currentBlock.Type = "code"
				currentBlock.Content = line + "\n"
				continue
			}
		}

		if inCodeBlock {
			currentBlock.Content += line + "\n"
			continue
		}

		if strings.HasPrefix(line, "#") {
			if currentBlock.Content != "" {
				blocks = append(blocks, currentBlock)
				currentBlock = Block{}
			}

			level := 0
			for _, ch := range line {
				if ch == '#' {
					level++
				} else {
					break
				}
			}

			headerContent := strings.TrimSpace(line[level:])
			tags := extractTags(headerContent)
			cleanContent := tagPattern.ReplaceAllString(headerContent, "")

			currentBlock.Type = "header"
			currentBlock.Level = level
			currentBlock.Content = strings.TrimSpace(cleanContent)
			currentBlock.Tags = tags
			blocks = append(blocks, currentBlock)
			currentBlock = Block{}
			continue
		}

		if strings.TrimSpace(line) == "" {
			if currentBlock.Content != "" {
				blocks = append(blocks, currentBlock)
				currentBlock = Block{}
			}
			continue
		}

		if currentBlock.Content == "" {
			currentBlock.Type = "text"
		}
		currentBlock.Content += line + "\n"
	}

	if currentBlock.Content != "" {
		blocks = append(blocks, currentBlock)
	}

	return &Module{Blocks: blocks}, nil
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
