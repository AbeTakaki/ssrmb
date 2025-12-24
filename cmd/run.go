package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "整理と削除を実行します",
	Run: func(cmd *cobra.Command, args []string) {
		dir := findScreenshotDir()
		if dir == "" {
			fmt.Println("❌ スクリーンショットフォルダが見つかりませんでした。")
			return
		}

		fmt.Printf("📂 対象: %s\n", dir)

		// 1. 整理の実行
		organize(dir)

		// 2. 削除の実行
		now := time.Now()
		if rmm > 0 {
			confirmAndClean(dir, now.AddDate(0, -rmm, 0))
		} else if rmd > 0 {
			confirmAndClean(dir, now.AddDate(0, 0, -rmd))
		}

		fmt.Println("\n✨ すべての処理が完了しました。")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// --- 以下、ロジック関数 ---

func findScreenshotDir() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "Pictures", "Screenshots"),
		filepath.Join(home, "OneDrive", "画像", "Screenshots"),
		filepath.Join(home, "OneDrive", "Pictures", "Screenshots"),
		filepath.Join(home, "画像", "スクリーンショット"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func organize(baseDir string) {
	fmt.Println("🧹 画像を日付フォルダに整理中...")
	files, _ := os.ReadDir(baseDir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		info, _ := f.Info()
		dateDir := info.ModTime().Format("2006-01-02")
		targetDir := filepath.Join(baseDir, dateDir)

		os.MkdirAll(targetDir, 0755)
		src := filepath.Join(baseDir, f.Name())
		dst := getSafePath(targetDir, f.Name())
		os.Rename(src, dst)
	}
}

func getSafePath(dir, filename string) string {
	path := filepath.Join(dir, filename)
	base := filename[:len(filename)-len(filepath.Ext(filename))]
	ext := filepath.Ext(filename)
	for i := 1; ; i++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}
		path = filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
	}
}

func confirmAndClean(baseDir string, threshold time.Time) {
	files, _ := os.ReadDir(baseDir)
	var targets []string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		date, err := time.Parse("2006-01-02", f.Name())
		if err == nil && date.Before(threshold) {
			targets = append(targets, f.Name())
		}
	}

	if len(targets) == 0 {
		return
	}

	fmt.Printf("\n⚠️  以下の古いフォルダを完全に削除します (基準: %s 以前):\n", threshold.Format("2006-01-02"))
	for _, t := range targets {
		fmt.Printf("  - %s\n", t)
	}
	fmt.Print("本当によろしいですか？ (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(input)) == "y" {
		for _, t := range targets {
			os.RemoveAll(filepath.Join(baseDir, t))
			fmt.Printf("🗑️  削除しました: %s\n", t)
		}
	} else {
		fmt.Println("🚫 削除をキャンセルしました。")
	}
}
