package configure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ConfigureFile = "go-daemon.json"
)

var g_Configure Configure

type Configure struct {
	//服務名稱
	Name string
	//服務 顯示名稱
	Show string
	//服務 描述
	Description string
	//啓動方式是否爲 自動
	Auto bool

	//被守衛進程 路徑
	Bin string
	//被守衛進程 工作路徑
	Directory string
	//被守衛進程 執行參數
	Params string
}

func (c *Configure) String() string {
	b, _ := json.MarshalIndent(&g_Configure, "", "	")
	return string(b)
}

func Init() error {
	file := filepath.Dir(os.Args[0]) + "/" + ConfigureFile

	b, e := ioutil.ReadFile(file)
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, &g_Configure)
	if e != nil {
		return e
	}
	return nil
}
func GetConfigure() *Configure {
	return &g_Configure
}
