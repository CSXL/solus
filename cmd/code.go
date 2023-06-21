package cmd

import (
	"fmt"

	"github.com/CSXL/solus/code"
	"github.com/spf13/cobra"
)

var GenerationFolder string

func init() {
	codeCmd.PersistentFlags().StringVarP(&GenerationFolder, "generation-folder", "g", "", "The folder to generate code in.")
	_ = codeCmd.MarkFlagRequired("generation-folder")
	rootCmd.AddCommand(codeCmd)
}

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Generate the code for a project",
	Long:  `Generate the code for a project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating code...")
		codeConfig, err := code.LoadCodeConfig(GenerationFolder)
		if err != nil {
			fmt.Println(err)
			return
		}
		codeGenerator := code.NewCodeGenerator(GenerationFolder, codeConfig)
		err = codeGenerator.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Generated code successfully!")
	},
}
