package cmd

import (
	"fmt"

	"github.com/CSXL/solus/requirements"
	"github.com/spf13/cobra"
)

var PathToConversation string
var OutputFile string

func init() {
	requirementsCmd.PersistentFlags().StringVarP(&PathToConversation, "conversation-file", "f", "", "The path to the conversation file.")
	_ = requirementsCmd.MarkFlagRequired("conversation-file")
	requirementsCmd.PersistentFlags().StringVarP(&OutputFile, "output-file", "o", "", "Write the generated requirements to a file.")
	rootCmd.AddCommand(requirementsCmd)
}

var requirementsCmd = &cobra.Command{
	Use:   "requirements",
	Short: "Generate requirements for a project from a file",
	Long:  `Generate requirements for a project from a file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating requirements...")
		requirementsConfig, err := requirements.LoadRequirementsConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
		inputConversation, err := requirements.LoadInputConversation(PathToConversation)
		if err != nil {
			fmt.Println(err)
			return
		}
		requirementsGenerator := requirements.NewRequirementsGenerator(inputConversation, requirementsConfig)
		_, err = requirementsGenerator.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Generated requirements successfully!")
		if OutputFile == "" {
			fmt.Println("Generated requirements: ", requirementsGenerator.GeneratedRequirements)
		} else {
			err = requirementsGenerator.SaveGeneratedRequirements(OutputFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Saved generated requirements to " + OutputFile)
		}
	},
}
