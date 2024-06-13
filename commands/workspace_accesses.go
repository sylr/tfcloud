package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	flag "github.com/spf13/pflag"
)

type TeamAccessListCommand struct {
	Command
	organization string
}

func (c *TeamAccessListCommand) Run(args []string) int {
	var formatOpt string
	f := flag.NewFlagSet("team_accesses", flag.ExitOnError)
	f.StringVarP(&formatOpt, "format", "f", "", "Output format. Available formats: json, table")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Error(err.Error())
		return 1
	}

	if formatOpt != "" {
		c.Command.Format = Format(formatOpt)
	}

	c.organization = f.Arg(0)

	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client

	aslist, err := c.Client.TeamAccessList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case FormatAlfred:
		alfredItems := make([]AlfredFormatItem, len(aslist))
		for i, v := range aslist {
			alfredItems[i] = AlfredFormatItem{
				Title:        v.Team.Name,
				SubTitle:     "",
				Arg:          fmt.Sprintf("%s/app/%s/workspaces/%s/variables/setting/access/%s", c.Client.Address(), c.organization, v.Workspace.Name, v.ID),
				Match:        v.Team.ID,
				AutoComplete: v.Team.Name,
				UID:          v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No workspaces found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	case FormatJSON:
		out, err := json.MarshalIndent(aslist, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "TEAM\tTEAM ID\tWORKSPACE\tACCESS ID\tPRIVILEGE")
		for _, v := range aslist {
			if v != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", v.Team.Name, v.Team.ID, v.Workspace.Name, v.ID, v.Access)
			}
		}
		w.Flush()
		c.UI.Output(out.String())
	}
	return 0
}

func (c *TeamAccessListCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceAccesses)
}

func (c *TeamAccessListCommand) Synopsis() string {
	return "Lists all terraform cloud workspaces"
}

const helpWorkspaceAccesses = `
Usage: tfcloud workspace list [OPTIONS] <organization>

  Lists all terraform cloud workspaces

Options:
  --format, -f             Output format. Available formats: json, table (default: table)
`
