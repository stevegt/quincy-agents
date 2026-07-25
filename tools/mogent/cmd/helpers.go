package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func findModuleFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	sort.Strings(files)
	return files
}

func extractTagsFromFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var tags []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if idx := strings.Index(line, "{#"); idx >= 0 {
			rest := line[idx+2:]
			if endIdx := strings.Index(rest, "}"); endIdx >= 0 {
				tagPart := rest[:endIdx]
				for _, t := range strings.Fields(tagPart) {
					t = strings.TrimPrefix(t, "#")
					if t != "" {
						tags = append(tags, t)
					}
				}
			}
		}
	}
	return tags
}

func printTagTree(tags map[string]bool) {
	roots := make(map[string][]string)
	for tag := range tags {
		parts := strings.SplitN(tag, "/", 2)
		if len(parts) == 1 {
			roots[tag] = nil
		} else {
			roots[parts[0]] = append(roots[parts[0]], parts[1])
		}
	}

	sortedRoots := make([]string, 0, len(roots))
	for r := range roots {
		sortedRoots = append(sortedRoots, r)
	}
	sort.Strings(sortedRoots)

	for i, root := range sortedRoots {
		children := roots[root]
		isLast := i == len(sortedRoots)-1
		prefix := "  "
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if len(children) == 0 {
			fmt.Printf("%s%s%s\n", prefix, connector, root)
		} else {
			fmt.Printf("%s%s%s/\n", prefix, connector, root)
			sort.Strings(children)
			for j, child := range children {
				childPrefix := prefix + "│   "
				if isLast {
					childPrefix = prefix + "    "
				}
				childConnector := "├── "
				if j == len(children)-1 {
					childConnector = "└── "
				}
				fmt.Printf("%s%s%s\n", childPrefix, childConnector, child)
			}
		}
	}
}
