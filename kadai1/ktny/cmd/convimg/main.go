package main

import (
	"flag"
	"fmt"
	"kadai1/ktny/converter"
	"os"
	"path/filepath"
	"strings"
)

var (
	from string
	to   string
)

func init() {
	flag.StringVar(&from, "from", converter.JPG, "from ext. support jpg, jpeg, png, gif.")
	flag.StringVar(&to, "to", converter.PNG, "to ext. support jpg, jpeg, png, gif.")
}

func main() {
	flag.Parse()

	// ディレクトリが指定されていない場合は終了する
	targetDir := flag.Arg(0)
	if targetDir == "" {
		fmt.Println("[Error]Directory is not defined.")
		os.Exit(1)
	}
	targetDir = filepath.Join(strings.Split(targetDir, "/")...)

	// 指定された変換前後の拡張子が同じ場合は終了する
	if from == to {
		fmt.Println("[Error]from and to extension is same.")
		os.Exit(1)
	}

	// 指定された拡張子がサポートされていない場合は終了する
	if !converter.IsSupportedExt(from) {
		fmt.Printf("[Error]%s is not supported ext.", from)
	} else if !converter.IsSupportedExt(to) {
		fmt.Printf("[Error]%s is not supported ext.", to)
	}

	fmt.Printf("[Info]from=%s, to=%s, targetDir=%s\n", from, to, targetDir)

	// targetDir配下のファイルパスを取得する
	filepaths, err := converter.DirWalk(targetDir)
	if err != nil {
		fmt.Printf("[Error]%s", err)
		os.Exit(1)
	}

	// 指定された拡張子の画像ファイルを変換する
	for _, filepath := range filepaths {
		if converter.CanConvert(from, filepath) {
			fmt.Printf("[Info]convert %s to %s\n", filepath, to)
			if err := converter.ConvertImage(filepath, from, to); err != nil {
				fmt.Printf("[Error]%s", err)
			}
		} else {
			fmt.Printf("[Warn]cannnot convert %s. It is not %s file.\n", filepath, from)
		}
	}
}
