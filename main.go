// main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	const maxLen = 32767
	const prefix = `\\?` + "\\"
	const fileName = "myfile.txt"

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Getwd error:", err)
		os.Exit(1)
	}

	baseLen := len(prefix) + len(cwd) + len(fileName)

	T := maxLen - baseLen - 1
	if T < 1 {
		fmt.Fprintln(os.Stderr, "기본 경로가 너무 깁니다:", baseLen)
		os.Exit(1)
	}

	segCount := (T + 254) / 255
	//    sumLens = T - segCount
	sumLens := T - segCount

	baseSeg := sumLens / segCount
	extra := sumLens % segCount

	parts := []string{cwd}
	for i := 0; i < segCount; i++ {
		segLen := baseSeg
		if i < extra {
			segLen++
		}
		parts = append(parts, strings.Repeat("a", segLen))
	}
	parts = append(parts, fileName)

	longPath := prefix + filepath.Join(parts...)
	fmt.Printf("생성할 경로 길이: %d\n", len(longPath))

	if err := os.MkdirAll(filepath.Dir(longPath), os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, "MkdirAll error:", err)
		os.Exit(1)
	}

	f, err := os.Create(longPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Create error:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString("Long path (32767) test\n"); err != nil {
		fmt.Fprintln(os.Stderr, "WriteString error:", err)
		os.Exit(1)
	}

	fmt.Println("파일 생성 완료!")
}
