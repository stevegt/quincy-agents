package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [module]",
	Short: "Edit a module or the assembled AGENTS.md",
	Long:  "Open a module file or the assembled AGENTS.md in $EDITOR",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		var target string
		if len(args) == 0 {
			target = "AGENTS.md"
		} else {
			target = fmt.Sprintf(".mogent/modules/%s.md", args[0])
		}

		c := exec.Command(editor, target)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}
