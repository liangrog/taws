// Ansible commands and options
package cmd

import (
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

	ParamFilterName = "filter-name"
	DescFilterName  = "Filter result by given name string (full/partial)"

	ParamInventoryFile = "inventory-file"
	DescInventoryFile  = "Full file path for alternative inventory file"
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
	cmd.Flags().StringVar(&ops.AsgNameFilter, ParamFilterName, ops.AsgNameFilter, DescFilterName)
	cmd.Flags().StringVar(&ops.InventoryFile, ParamInventoryFile, ops.InventoryFile, DescInventoryFile)

	return cmd
}
