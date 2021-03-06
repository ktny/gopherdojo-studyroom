// 画像変換パッケージ
package converter

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type img image.Image

// サポートする拡張子
const (
	JPG  = "jpg"
	JPEG = "jpeg"
	PNG  = "png"
	GIF  = "gif"
)

// extがサポートされている拡張子かを返す
func IsSupportedExt(ext string) bool {
	switch ext {
	case JPG, JPEG, PNG, GIF:
		return true
	default:
		return false
	}
}

// dir配下のファイルパスを再帰的に返す
func DirWalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			childPaths, err := DirWalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, childPaths...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
}

// 変換可能なファイルかを返す
func CanConvert(from, path string) bool {
	switch from {
	case JPG, JPEG:
		return strings.ToLower(filepath.Ext(path)) == toExt(JPG) || strings.ToLower(filepath.Ext(path)) == toExt(JPEG)
	case PNG:
		return strings.ToLower(filepath.Ext(path)) == toExt(PNG)
	case GIF:
		return strings.ToLower(filepath.Ext(path)) == toExt(GIF)
	default:
		return false
	}
}

// 文字列の先頭に"."を付加して返す
func toExt(s string) string {
	return "." + s
}

// 画像を指定の拡張子で変換する
func ConvertImage(path, from, to string) error {
	buf := new(bytes.Buffer)
	newFilePath := strings.Replace(path, from, to, 1)

	img, err := decodeImg(path, from)
	if err != nil {
		return err
	}

	if err := encodeImg(buf, &img, to); err != nil {
		return err
	}

	if err := updateImg(buf, path, newFilePath); err != nil {
		return err
	}

	return nil
}

// 画像をデコードする
func decodeImg(path, from string) (img, error) {
	var img img

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch from {
	case JPG, JPEG:
		img, err = jpeg.Decode(f)
		if err != nil {
			return nil, err
		}
	case PNG:
		img, err = png.Decode(f)
		if err != nil {
			return nil, err
		}
	case GIF:
		img, err = gif.Decode(f)
		if err != nil {
			return nil, err
		}
	}

	return img, err
}

// 画像をtoでエンコードする
func encodeImg(buf *bytes.Buffer, img *img, to string) error {
	switch to {
	case JPG, JPEG:
		options := &jpeg.Options{Quality: 100}
		if err := jpeg.Encode(buf, *img, options); err != nil {
			return err
		}
	case PNG:
		if err := png.Encode(buf, *img); err != nil {
			return err
		}
	case GIF:
		options := &gif.Options{}
		if err := gif.Encode(buf, *img, options); err != nil {
			return err
		}
	}
	return nil
}

// 画像を削除して、エンコードした画像で再作成する
func updateImg(buf *bytes.Buffer, path, newFilePath string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	file, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	var rerr error
	defer func() {
		if err := file.Close(); err != nil {
			rerr = err
		}
	}()
	if rerr != nil {
		return err
	}

	file.Write(buf.Bytes())
	return nil
}
