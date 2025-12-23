package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	rmm := flag.Int("rmm", 0, "指定した月数より前のフォルダを削除")
	rmd := flag.Int("rmd", 0, "指定した日数より前のフォルダを削除")
	flag.Parse()

	screenshotDir := findScreenshotDir()
	if screenshotDir == "" {
		fmt.Println("エラー: スクリーンショットフォルダが見つかりませんでした。")
		return
	}

	fmt.Printf("対象ディレクトリ: %s\n", screenshotDir)

	// 1. 画像の整理（移動）
	organizeScreenshots(screenshotDir)

	// 2. 削除処理の準備
	now := time.Now()
	var threshold time.Time
	if *rmm > 0 {
		threshold = now.AddDate(0, -*rmm, 0)
		confirmAndClean(screenshotDir, threshold)
	} else if *rmd > 0 {
		threshold = now.AddDate(0, 0, -*rmd)
		confirmAndClean(screenshotDir, threshold)
	}

	fmt.Println("処理が完了しました。")
}

// フォルダを探索する関数（前回と同じ）
func findScreenshotDir() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "Pictures", "Screenshots"),
		filepath.Join(home, "OneDrive", "画像", "Screenshots"),
		filepath.Join(home, "OneDrive", "Pictures", "Screenshots"),
		filepath.Join(home, "画像", "スクリーンショット"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// 整理ロジック（前回と同じ）
func organizeScreenshots(baseDir string) {
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

		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			os.MkdirAll(targetDir, 0755)
		}

		srcPath := filepath.Join(baseDir, f.Name())
		destPath := getSafePath(targetDir, f.Name())
		os.Rename(srcPath, destPath)
	}
}

func getSafePath(dir, filename string) string {
	base := filename[:len(filename)-len(filepath.Ext(filename))]
	ext := filepath.Ext(filename)
	path := filepath.Join(dir, filename)
	for i := 1; ; i++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}
		path = filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
	}
}

// --- 今回新しく追加・修正した削除確認ロジック ---

func confirmAndClean(baseDir string, threshold time.Time) {
	files, _ := os.ReadDir(baseDir)
	var targets []string

	fmt.Printf("\n--- 削除確認 (基準日: %s 以前) ---\n", threshold.Format("2006-01-02"))

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		folderDate, err := time.Parse("2006-01-02", f.Name())
		if err != nil {
			continue
		}

		if folderDate.Before(threshold) {
			targets = append(targets, f.Name())
		}
	}

	if len(targets) == 0 {
		fmt.Println("削除対象の古いフォルダは見つかりませんでした。")
		return
	}

	fmt.Println("以下のフォルダを完全に削除します:")
	for _, name := range targets {
		fmt.Printf("  [DIR] %s\n", name)
	}

	fmt.Print("\n本当によろしいですか？ (y/n): ")

	// ユーザー入力を受け取る
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		for _, name := range targets {
			err := os.RemoveAll(filepath.Join(baseDir, name))
			if err != nil {
				fmt.Printf("削除失敗: %s (%v)\n", name, err)
			} else {
				fmt.Printf("削除済み: %s\n", name)
			}
		}
	} else {
		fmt.Println("削除をキャンセルしました。")
	}
}
