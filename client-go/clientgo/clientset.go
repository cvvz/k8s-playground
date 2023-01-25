package main

import (
	"fmt"

	"github.com/pingcap/tidb-operator/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/chenweizhi/.kube/config")
	if err != nil {
		panic(err)
	}
	/************ kubernetes *************/
	// clientset := kubernetes.NewForConfigOrDie(config)

	// podClient := clientset.CoreV1().Pods("blade")

	// pod, err := podClient.Get("inf-blade-bladedevbasic-tikv-0", metav1.GetOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// // memoryInBytes := pod.Spec.Containers[0].Resources.Requests.Memory().Value()
	// // memoryTotal := float32(memoryInBytes/1048576) * 0.75
	// storageQuantity := pod.Spec.Containers[0].Resources.Requests[corev1.ResourceStorage]
	// storageInBytes := storageQuantity.Value()
	// fmt.Printf("%d\n", storageInBytes)

	/************ tidb cluster *************/
	cli := versioned.NewForConfigOrDie(config)

	tcClient := cli.PingcapV1alpha1().TidbClusters("blade")

	tc, err := tcClient.Get("inf-blade-bladedevbasic", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	storageQuantity := tc.Spec.TiKV.Requests[corev1.ResourceStorage]
	storageInBytes := storageQuantity.Value()
	storageInTerabyte := float64(storageQuantity.Value()) / 1099511627776
	fmt.Printf("%d bytes, %f TB\n", storageInBytes, storageInTerabyte)

}
