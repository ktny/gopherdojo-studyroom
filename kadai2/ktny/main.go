package main

import (
	"log"
	"os"
)

func readSample() {
	m := 5
	filename := "test.txt" // ABC
	f, err := os.Open(filename)
	defer f.Close()
	p := make([]byte, m)
	n, err := f.Read(p)
	if m < n {
		log.Fatalf("%dバイト読もうとしましたが、%dバイトしか読めませんでした\n", n, m)
	}
	if err != nil {
		log.Fatalf("読み込み中にエラーが発生しました:%v\n。", err)
	}
	log.Printf("読み込んだbyte数: %v", n) // 3
	log.Printf("中身: %v", p)         // [65 66 67 0 0]
}

func writeSample() {
	filename := "test.txt"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := "ABC"
	n, err := f.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	println(n) // 3
}

func main() {
	writeSample()
	readSample()
}
