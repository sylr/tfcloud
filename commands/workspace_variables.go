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

type WorkspaceVariablesCommand struct {
	Command
	organization string
}

func (c *WorkspaceVariablesCommand) Run(args []string) int {
	var formatOpt string
	f := flag.NewFlagSet("workspace_variables", flag.ExitOnError)
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

	vslist, err := c.Client.WorkspaceVariableList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case FormatAlfred:
		alfredItems := make([]AlfredFormatItem, len(vslist))
		for i, v := range vslist {
			alfredItems[i] = AlfredFormatItem{
				Title:        v.Variable.ID,
				SubTitle:     fmt.Sprintf("vcs repo: %s", *v.Workspace.VCSRepoName),
				Arg:          fmt.Sprintf("%s/app/%s/workspaces/%s/variables/%s", c.Client.Address(), c.organization, *v.Workspace.Name, v.Variable.ID),
				Match:        v.Variable.Key,
				AutoComplete: v.Variable.Value,
				UID:          v.Variable.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No workspaces found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	case FormatJSON:
		out, err := json.MarshalIndent(vslist, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "WORKSPACE\tNAME\tID\tVALUE\t")
		for _, v := range vslist {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", *v.Workspace.Name, v.Variable.Key, v.Variable.ID, v.Variable.Value)
		}
		w.Flush()
		c.UI.Output(out.String())
	}
	return 0
}

func (c *WorkspaceVariablesCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceVariables)
}

func (c *WorkspaceVariablesCommand) Synopsis() string {
	return "Lists all terraform cloud workspaces"
}

const helpWorkspaceVariables = `
Usage: tfcloud workspace list [OPTIONS] <organization>

  Lists all terraform cloud workspaces

Options:
  --format, -f             Output format. Available formats: json, table (default: table)
`
