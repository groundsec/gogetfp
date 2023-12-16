package gogetfp

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FreeProxy struct defines the parameters for scraping and checking proxies.
type FreeProxyConfig struct {
	CountryID []string
	Timeout   float64
	Random    bool
	Anonym    bool
	Elite     bool
	Google    *bool
	HTTPS     bool
	Schema    string
}

type FreeProxy struct {
	Config FreeProxyConfig
}

var DefaultFreeProxyConfig = FreeProxyConfig{
	CountryID: []string{},
	Timeout:   0.5,
	Random:    false,
	Anonym:    false,
	Elite:     false,
	Google:    nil,
	HTTPS:     false,
	Schema:    "http",
}

func New(config FreeProxyConfig) *FreeProxy {
	if config.Timeout == 0 {
		config.Timeout = DefaultFreeProxyConfig.Timeout
	}
	if config.Schema == "" {
		config.Schema = DefaultFreeProxyConfig.Schema
	}
	return &FreeProxy{Config: config}
}

// GetProxyList retrieves a list of proxies based on the specified criteria.
func (fp *FreeProxy) GetProxyList(repeat bool) ([]string, error) {
	var website string
	if repeat {
		website = "https://free-proxy-list.net"
	} else if slices.Contains(fp.Config.CountryID, "US") {
		website = "https://www.us-proxy.org"
	} else if slices.Contains(fp.Config.CountryID, "GB") {
		website = "https://free-proxy-list.net/uk-proxy.html"
	} else {
		website = "https://www.sslproxies.org"
	}

	resp, err := http.Get(website)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %v", website, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %v", err)
	}

	var proxies []string
	doc.Find("#list tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			row := s.Children()
			if fp.criteria(row) {
				proxy := fmt.Sprintf("%s:%s", row.Eq(0).Text(), row.Eq(1).Text())
				proxies = append(proxies, proxy)
			}
		}
	})
	return proxies, nil
}

// criteria checks if a given row of elements meets the specified criteria.
func (fp *FreeProxy) criteria(row *goquery.Selection) bool {
	countryCriteria := len(fp.Config.CountryID) == 0 || row.Eq(2).Text() == fp.Config.CountryID[0]
	eliteCriteria := len(fp.Config.CountryID) == 0 || strings.Contains(row.Eq(4).Text(), "elite")
	anonymCriteria := len(fp.Config.CountryID) == 0 || !fp.Config.Anonym || strings.Contains(row.Eq(4).Text(), "anonymous")
	googleCriteria := fp.Config.Google == nil || *fp.Config.Google == (row.Eq(5).Text() == "yes")
	httpsCriteria := !fp.Config.HTTPS || strings.ToLower(row.Eq(6).Text()) == "yes"
	return countryCriteria && eliteCriteria && anonymCriteria && googleCriteria && httpsCriteria
}

// Get returns a working proxy that matches the specified parameters.
func (fp *FreeProxy) Get(repeat bool) (string, error) {
	proxyList, err := fp.GetProxyList(repeat)
	if err != nil {
		return "", err
	}

	if fp.Config.Random {
		rand.Shuffle(len(proxyList), func(i, j int) {
			proxyList[i], proxyList[j] = proxyList[j], proxyList[i]
		})
	}

	var workingProxy string
	for _, proxyAddress := range proxyList {
		workingProxy, err = fp.checkIfProxyIsWorking(proxyAddress)
		if err == nil && workingProxy != "" {
			return workingProxy, nil
		}
	}

	if workingProxy == "" && !repeat {
		fp.Config.CountryID = nil
		return fp.Get(true)
	}

	return "", fmt.Errorf("there are no working proxies at this time")
}

// checkIfProxyIsWorking checks if a proxy is working by making a request to Google.
func (fp *FreeProxy) checkIfProxyIsWorking(proxyAddress string) (string, error) {
	testUrl := fmt.Sprintf("%s://www.google.com", fp.Config.Schema)
	proxy := fmt.Sprintf("%s://%s", fp.Config.Schema, proxyAddress)
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return "", err
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := client.Get(testUrl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	if conn, ok := resp.Body.(interface{ CloseIdleConnections() }); ok {
		conn.CloseIdleConnections()
	}

	if resp.Request != nil && resp.Request.URL != nil && resp.StatusCode == 200 {
		return proxy, nil
	}

	return "", nil
}
