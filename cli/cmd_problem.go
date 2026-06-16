package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tamnd/uva-cli/uva"
)

func (a *App) problemCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "problem <num>",
		Short: "Get a UVa problem by number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return codeError(exitUsage, fmt.Errorf("invalid problem number: %s", args[0]))
			}
			p, err := a.client.GetByNum(cmd.Context(), num)
			if err != nil {
				return mapFetchErr(err)
			}
			if p == nil {
				return codeError(exitNoData, fmt.Errorf("problem %d not found", num))
			}
			return a.render([]uva.Problem{*p})
		},
	}
}
