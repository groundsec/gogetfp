package gogetfp

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

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
}

type FreeProxy struct {
	Config FreeProxyConfig
}

var defaultFreeProxyConfig = FreeProxyConfig{
	CountryID: []string{},
	Timeout:   5,
	Random:    false,
	Anonym:    false,
	Elite:     false,
	Google:    nil,
	HTTPS:     false,
}

func New(config FreeProxyConfig) *FreeProxy {
	if config.Timeout == 0 {
		config.Timeout = defaultFreeProxyConfig.Timeout
	}
	return &FreeProxy{Config: config}
}

// GetProxyList retrieves a list of proxies based on the specified criteria.
func (fp *FreeProxy) GetProxyList() ([]string, error) {
	var website string
	if slices.Contains(fp.Config.CountryID, "US") {
		website = "https://www.us-proxy.org"
	} else if slices.Contains(fp.Config.CountryID, "GB") {
		website = "https://free-proxy-list.net/uk-proxy.html"
	} else {
		defaultProxies := []string{"https://www.sslproxies.org", "https://free-proxy-list.net"}
		randomIndex := rand.Intn(2)
		website = defaultProxies[randomIndex]
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

func (fp *FreeProxy) getSchema() string {
	if fp.Config.HTTPS {
		return "https"
	} else {
		return "http"
	}
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
func (fp *FreeProxy) Get() (string, error) {
	proxyList, err := fp.GetProxyList()
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

	return "", fmt.Errorf("there are no working proxies at this time")
}

// checkIfProxyIsWorking checks if a proxy is working by making a request to Google.
func (fp *FreeProxy) checkIfProxyIsWorking(proxyAddress string) (string, error) {
	schema := fp.getSchema()
	testUrl := fmt.Sprintf("%s://www.example.com", schema)
	proxy := fmt.Sprintf("%s://%s", schema, proxyAddress)
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return "", err
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}, Timeout: time.Duration(fp.Config.Timeout) * time.Second}
	resp, err := client.Get(testUrl)
	if err != nil {
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
