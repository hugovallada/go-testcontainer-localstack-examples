package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("AWS ENDPOINT %v\n", os.Getenv("AWS_ENDPOINT"))
}
