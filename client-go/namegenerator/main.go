package main

import (
	"fmt"

	"k8s.io/apiserver/pkg/storage/names"
)

func main() {
	podName := "test-pod-"
	fmt.Println(names.SimpleNameGenerator.GenerateName(podName))
}
