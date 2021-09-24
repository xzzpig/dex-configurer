package utils_test

import (
	"strings"
	"testing"

	"github.com/xzzpig/dex-configurer/utils"
)

func TestMd5(t *testing.T) {
	t.Log(utils.Md5("MD5testing") == "f7bb96d1dcd6cfe0e5ce1f03e35f84bf")
}

func TestMarshalCamel(t *testing.T) {
	str := "de-sample"
	str = strings.ReplaceAll(str, "-", "_")
	t.Log(str)
	str = utils.MarshalCamel(str)
	t.Log(str)
	str = utils.UnMarshalCamel(str)
	t.Log(str)
}

func TestRandString(t *testing.T) {
	t.Log(utils.RandString(8))
}
