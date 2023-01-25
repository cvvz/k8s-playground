package main

import (
	"bytes"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/chenweizhi/.kube/config")
	if err != nil {
		panic(err)
	}

	k8sCli, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	cmd := []string{
		"sh",
		"-c",
		"df -BG /mnt/azuredisk | awk 'NR == 2'  | awk {'print $2'}",
	}

	req := k8sCli.CoreV1().RESTClient().Post().
		Resource("pods").
		Name("statefulset-azuredisk-0").
		Namespace("default").SubResource("exec")
	req.VersionedParams(
		&v1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		},
		scheme.ParameterCodec,
	)

	var stdout, stderr bytes.Buffer
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.EqualFold(strings.TrimSpace(stdout.String()), "15G"))
	return
}
