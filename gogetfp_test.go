package gogetfp_test

import (
	"fmt"
	"testing"

	"github.com/groundsec/gogetfp"
)

func TestDefaultProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{})

	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestRandomProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{Random: true})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestAnonProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{Anonym: true})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestEliteProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{Elite: true})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestHTTPSProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{HTTPS: true})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestGoogleProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{HTTPS: true})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}

func TestGBProxy(t *testing.T) {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{Timeout: 5, CountryID: []string{"GB"}})
	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}
