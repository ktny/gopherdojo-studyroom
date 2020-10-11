# Gopher道場 課題2
## io.Readerとio.Writerについて調べる
### 標準パッケージでどのように使われているか

#### io.Readerのインターフェース定義
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
[io.Reader](https://golang.org/pkg/io/#Reader)

読み込んだデータをpにbyteのスライスとして埋め、読み込んだbyte数をnで返すという実装になっている。
ここで読み込む先のインターフェースの種類は規定されておらず、ファイル、メモリ、ネットワークなど読み込み先はio.Readerを実装する型によって規定される仕組みになっている。

例えば、ファイルを読み込む場合はos.File*型にReadメソッドが用意されており、io.Readerのシグネチャと同一なので、os.File*型はio.Readerインタフェースを実装していると言える。
```go
func (f *File) Read(b []byte) (n int, err error)
```
[os.File*.Readメソッド](https://golang.org/pkg/os/#File.Read)

実際にos.File*型でReadメソッドを使用すると次のようになる。
```go
func main() {
	m := 5
	filename := "test.txt" // ABC
	f, err := os.Open(filename)
	p := make([]byte, m)
    n, err := f.Read(p)
	if m < n {
		log.Fatalf("%dバイト読もうとしましたが、%dバイトしか読めませんでした\n", n, m)
	}
	if err != nil {
		log.Fatalf("読み込み中にエラーが発生しました:%v\n。", err)
	}
	log.Printf("読み込んだbyte数: %v", n) // 3
	log.Printf("中身: %v", p) // [65 66 67 0 0]
}
```

また、他の例では、bufio.NewReaderメソッドではio.Readerインターフェースを実装する型を引数に、bufio.Reader型に包んで返すことができる。
bufio.Reader型はio.Reader型にバッファリングの機能を実装するもので、通常であればReadメソッドが呼び出される度にシステムコールが発生するのをデフォルトで4096バイトごとに読み込みを行うようバッファするのでシステムコールの回数が減り負荷を減らすことができる。
```go
func NewReader(rd io.Reader) *Reader
```
[bufio.NewReaderメソッド](https://golang.org/pkg/bufio/#NewReader)

#### io.Writerのインターフェース定義
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
[io.Writer](https://golang.org/pkg/io/#Writer)

byteのスライス型であるpを引数に渡し、pを書き込んだbyte数をnで返すという実装になっている。Readと同じシグネチャ。
こちらもReadと同じように書き込む先のインターフェースはio.Writerを実装する型によって規定される仕組みになっている。

例えば、ファイルに書き込む場合はReadと同じくos.File*型がio.Writerを実装している。
```go
func (f *File) Write(b []byte) (n int, err error)
```
[os.File*.Writeメソッド](https://golang.org/pkg/os/#File.Write)

実際にos.File*型でWriteメソッドを使用すると次のようになる。
```go
func main() {
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
```
書き込みたい文字列を`[]byte`型でキャストすることにより、Writeメソッドに渡せるようになり指定ファイルに書き込めるようになる。

このようにgoのio.Reader, io.Writerインターフェースは標準パッケージでも、汎用的な入出力を行う実装としてよく使用されている。

### io.Readerとio.Writerがあることでどういう利点があるのか


