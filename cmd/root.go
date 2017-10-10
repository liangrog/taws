package cmd

import (
	"github.com/spf13/cobra"
)

const (
	ParamProfile = "profile"
	DescProfile  = "AWS CLI profile name"

	ParamRegion = "region"
	DescRegion  = "AWS region to access"
)

var profile string
var region string

var RootCmd = GetRootCmd()

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "taws",
		Short: "Tools for AWS",
		Long:  `Utility tooling for AWS`,
		/*Run: func(cmd *cobra.Command, args []string) {
		},*/
	}

	rootCmd.PersistentFlags().StringVar(&profile, ParamProfile, profile, DescProfile)
	rootCmd.PersistentFlags().StringVar(&region, ParamRegion, region, DescRegion)

	return rootCmd
}
