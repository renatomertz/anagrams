package main

import (
	"fmt"
	"time"

	"rmertz.com/anagram/internal/service"
)

func main() {
	start := time.Now()

	in := []string{"H", "I", "F", "A", "A", "L"} // "O"}
	service.ProcessParallelVersion(in)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
