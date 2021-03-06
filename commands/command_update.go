package commands

import (
	"regexp"

	"github.com/github/git-lfs/git"
	"github.com/github/git-lfs/lfs"
	"github.com/spf13/cobra"
)

var (
	updateForce  = false
	updateManual = false
)

// updateCommand is used for updating parts of Git LFS that reside under
// .git/lfs.
func updateCommand(cmd *cobra.Command, args []string) {
	requireInRepo()

	lfsAccessRE := regexp.MustCompile(`\Alfs\.(.*)\.access\z`)
	for key, value := range cfg.AllGitConfig() {
		matches := lfsAccessRE.FindStringSubmatch(key)
		if len(matches) < 2 {
			continue
		}

		switch value {
		case "basic":
		case "private":
			git.Config.SetLocal("", key, "basic")
			Print("Updated %s access from %s to %s.", matches[1], value, "basic")
		default:
			git.Config.UnsetLocalKey("", key)
			Print("Removed invalid %s access of %s.", matches[1], value)
		}
	}

	if updateForce && updateManual {
		Exit("You cannot use --force and --manual options together")
	}

	if updateManual {
		Print(lfs.GetHookInstallSteps())
	} else {
		if err := lfs.InstallHooks(updateForce); err != nil {
			Error(err.Error())
			Exit("To resolve this, either:\n  1: run `git lfs update --manual` for instructions on how to merge hooks.\n  2: run `git lfs update --force` to overwrite your hook.")
		} else {
			Print("Updated pre-push hook.")
		}
	}

}

func init() {
	RegisterSubcommand(func() *cobra.Command {
		cmd := &cobra.Command{
			Use: "update",
			Run: updateCommand,
		}
		cmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Overwrite existing hooks.")
		cmd.Flags().BoolVarP(&updateManual, "manual", "m", false, "Print instructions for manual install.")
		return cmd
	})
}
