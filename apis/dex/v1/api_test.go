package v1_test

import (
	"testing"

	v1 "github.com/xzzpig/dex-configurer/apis/dex/v1"
)

func TestDexRedirectUrlGoString(t *testing.T) {
	url := v1.DexRedirectUrl{
		Scheme: "https",
		Host:   "dex.example.com",
		Path:   "/oauth2",
		Port:   443,
	}
	t.Log(url.GoString())
	t.Log(url.AuthSignin())
	t.Log("=====================================")

	url = v1.DexRedirectUrl{
		Scheme: "https",
		Host:   "dex.example.com",
		Path:   "/oauth2",
		Port:   8443,
	}
	t.Log(url.GoString())
	t.Log(url.AuthSignin())
	t.Log("=====================================")

	url = v1.DexRedirectUrl{
		Scheme: "http",
		Host:   "dex.example.com",
		Path:   "/oauth2",
		Port:   80,
	}
	t.Log(url.GoString())
	t.Log(url.AuthSignin())
	t.Log("=====================================")

	url = v1.DexRedirectUrl{
		Scheme: "http",
		Host:   "dex.example.com",
		Path:   "/oauth2",
		Port:   8080,
	}
	t.Log(url.GoString())
	t.Log(url.AuthSignin())
	t.Log("=====================================")

	url = v1.DexRedirectUrl{
		Scheme: "http",
		Host:   "dex.example.com",
		Path:   "/oauth2/",
		Port:   8080,
	}
	t.Log(url.GoString())
	t.Log(url.AuthSignin())
}
