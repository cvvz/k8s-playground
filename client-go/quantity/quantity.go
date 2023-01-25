package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

func main() {
	storageSize := "10Gi"
	q := resource.MustParse(storageSize)
	fmt.Printf("%s\n %d\n %d\n", q.String(), q.Value(), q.ScaledValue(resource.Giga))
}
