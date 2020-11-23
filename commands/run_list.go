package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

type RunListCommand struct {
	UI cli.Ui
}

func (c *RunListCommand) Run(args []string) int {
	if len(args) < 3 {
		c.UI.Error("Arguments is not valid.")
		c.UI.Info(c.Help())
		return 1
	}

	buf := &bytes.Buffer{}
	var format string
	f := flag.NewFlagSet("run_list", flag.ContinueOnError)
	f.SetOutput(buf)
	f.StringVar(&format, "output", "table", "output format (table, json)")
	if err := f.Parse(args); err != nil {
		c.UI.Info(c.Help())
		return 1
	}
	if format != "table" && format != "json" {
		c.UI.Error("--output must be 'table' or 'json'")
		c.UI.Info(c.Help())
		return 1
	}

	organization := args[0]
	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	result, err := client.RunList(organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch format {
	case "table":
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "WORKSPACE\tSTATUS\tNEEDS CONFIRM\tLINK")
		for _, r := range result {
			fmt.Fprintf(w, "%s\t%s\t%v\thttps://%s/app/%s/workspaces/%s/runs/%s\n", r.Workspace, r.Status, r.IsConfirmable, address, organization, r.Workspace, r.ID)
		}
		w.Flush()
		c.UI.Output(out.String())
	case "json":
		out, err := json.Marshal(result)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	}
	return 0
}

func (c *RunListCommand) Help() string {
	return strings.TrimSpace(helpRunList)
}

func (c *RunListCommand) Synopsis() string {
	return "List all current terraform runs"
}

const helpRunList = `
Usage: tfcloud run list <organization>
`
