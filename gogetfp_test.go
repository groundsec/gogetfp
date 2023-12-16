package gogetfp_test

import (
	"fmt"
	"testing"

	"github.com/groundsec/gogetfp"
)

func TestDefaultProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{})

	proxy, err := fp.Get(false)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}
