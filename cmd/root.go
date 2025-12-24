/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// フラグ用の変数
var rmm int
var rmd int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssrmb",
	Short: "スクリーンショットの整理と削除を実行します。",
	Long:  `スクリーンショットフォルダにある画像を日付ごとに振り分けます。指定がある場合は古いフォルダを削除します。`,
	Run: func(cmd *cobra.Command, args []string) {
		runCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// グローバルフラグ
	rootCmd.PersistentFlags().IntVarP(&rmm, "rmm", "m", 0, "指定した月数より前のフォルダを削除")
	rootCmd.PersistentFlags().IntVarP(&rmd, "rmd", "d", 0, "指定した日数より前のフォルダを削除")
}
