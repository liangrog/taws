// Ansible commands and options
package cmd

import (
	"errors"
	_ "fmt"
	"github.com/liangrog/taws/ansible"
	"github.com/liangrog/taws/utils"
	"github.com/spf13/cobra"
)

// Defaults
const (
	// ansible
	CmdAnsible      = "ansible"
	CmdAnsibleShort = "Tools for ansible"
	CmdAnsibleLong  = `Tools for ansible`

	// get-inventory
	CmdGetInventory      = "get-inventory"
	CmdGetInventoryShort = "Get EC2 inventory"
	CmdGetInventoryLong  = `Get EC2 inventory`

	ParamToFile = "to-file"
	DescToFile  = "Full file path for alternative inventory file"

	ParamGroupBy = "group-by"
	DescGroupBy  = "Group result by. Available: asg"

	ParamFilterBy = "filter-by"
	DescFilterBy  = "Filter result by Available filters: tags, asg-name"

	ParamFilterValue = "filter-value"
	DescFilterValue  = "Filter values. if filtered by tags, use 'key1=value1;key2=value2' format. If filtered by asg-name, use 'name' string"

	ParamUsePublicIp = "use-public-ip"
	DescUsePublicIp  = "If to use public IP rather than private IP"
)

// Add ansible and sub command to root
func init() {
	AnsibleCmd := NewCmdAnsible()
	RootCmd.AddCommand(AnsibleCmd)

	AnsibleCmd.AddCommand(NewCmdGetInventory())
}

// Ansible command
func NewCmdAnsible() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdAnsible,
		Short: CmdAnsibleShort,
		Long:  CmdAnsibleLong,
	}

	return cmd
}

// Sub command "get-inventory"
func NewCmdGetInventory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdGetInventory,
		Short: CmdGetInventoryShort,
		Long:  CmdGetInventoryLong,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := filterValid(cmd); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Load AWS session from root persistent options
			utils.GetSession(
				cmd.InheritedFlags().Lookup("profile").Value.String(),
				cmd.InheritedFlags().Lookup("region").Value.String(),
			)

			ansible.GetInventory()
		},
	}

	// Add flags to command
	ops := ansible.NewOpsGetInventory()

	cmd.Flags().StringVar(&ops.GroupBy, ParamGroupBy, ops.GroupBy, DescGroupBy)
	cmd.Flags().StringVar(&ops.FilterBy, ParamFilterBy, ops.FilterBy, DescFilterBy)
	cmd.Flags().StringVar(&ops.FilterValue, ParamFilterValue, ops.FilterValue, DescFilterValue)
	cmd.Flags().BoolVar(&ops.UsePublicIp, ParamUsePublicIp, ops.UsePublicIp, DescUsePublicIp)
	cmd.Flags().StringVar(&ops.ToFile, ParamToFile, ops.ToFile, DescToFile)

	return cmd
}

// Check if filter options are valid
func filterValid(cmd *cobra.Command) error {
	fb := cmd.Flag(ParamFilterBy)
	t := cmd.Flag(ParamFilterValue)

	if len(fb.Value.String()) > 0 && len(t.Value.String()) == 0 {
		return errors.New("When you use --filter-by, flag --filter-value is required")
	}

	return nil
}
