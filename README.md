# gogetfp

`gogetfp` (**Go Get** **F**ree **P**roxy) is a lightweight Go library facilitating the retrieval of free proxies for seamless integration into your applications. With support for popular proxy sources like `free-proxy-list.net`, `sslproxies.org` and `us-proxy.org`. `gogetfp` offers customizable options, including **country filtering** and **timeout settings**. Enhance privacy, security, and access geo-restricted content effortlessly, making 'gogetfp' the go-to solution for integrating free proxies into your Go projects.

Heavily inspired by the Python library [freeproxy](https://github.com/jundymek/free-proxy).

## 🚀 Installation

Run the following command to install the latest version:

```bash
go install -v github.com/groundsec/gogetfp@latest
```

After this command `gogetfp` library source will be in your current `go.mod`.

## ⚙️ Options

The options are basically the same provided by [freeproxy](https://github.com/jundymek/free-proxy):

| Name      | Type     | Example      | Default value |
| --------- | -------- | ------------ | ------------- |
| CountryID | []string | ['US', 'BR'] | []string{}    |
| Timeout   | float64  | 0.1          | 1             |
| Random    | bool     | True         | False         |
| Anonym    | bool     | True         | False         |
| Elite     | bool     | True         | False         |
| Google    | bool     | False        | False         |
| HTTPS     | bool     | True         | False         |

## 💡 Example

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/groundsec/gogetfp"
)

func main() {
	fp := gogetfp.New(gogetfp.FreeProxyConfig{})

	proxy, err := fp.GetWorkingProxy()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Working Proxy:", proxy)
	}
}
```

## 🪪 License

_gogetfp_ is made with 🖤 by the [GroundSec](https://groundsec.io) team and released under the [MIT LICENSE](https://github.com/groundsec/gogetfp/blob/main/LICENSE).
