package main

import (
	"fmt"
	"os"

	cmd "github.com/bryant-rh/certcli/cmd/certcli"
	"github.com/bryant-rh/certcli/pkg/global"
	"k8s.io/klog/v2"
)

func main() {
	defer klog.Flush()
	cmd := cmd.NewCmd(&global.CertcliOptions{})
	if err := cmd.Execute(); err != nil {
		if klog.V(1).Enabled() {
			klog.Fatalf("%+v", err) // with stack trace
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
