package rules

import (
	"github.com/chrisenytc/skynet/config"
)

func Load(method string, path string) (string, []string) {

	if method == "GET" && path == "old" {
		return config.Get().OldApiUrl + "/?t=The+Flash", []string{}
	}

	if method == "GET" && path == "new" {
		return config.Get().NewApiUrl + "/3/discover/tv", []string{}
	}

	return "not_found", []string{}
}
