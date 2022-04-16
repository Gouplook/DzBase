package lang

import (
	"github.com/Gouplook/dzgin"
	"github.com/Gouplook/dzgin/i18n"
	"github.com/Gouplook/dzgin/utils"

	//"git.900sui.cn/kc/kcgin"
	//"git.900sui.cn/kc/kcgin/i18n"
	//"git.900sui.cn/kc/kcgin/utils"
	"log"
	"path/filepath"
)

//定义变量
var (
	loadLangs   map[string]*Load
	defaultLang string
)

//结构体
type Load struct {
	i18n *i18n.Locale
}

//获取lang
func GetLang(str string, langs ...string) string {
	for _, lang := range langs {
		if msg := getLang(str, lang); msg != "" {
			return msg
		}
	}
	if msg := getLang(str, defaultLang); msg != "" {
		return msg
	}
	return ""
}

//返回值
func getLang(str string, lang string) string {
	if loadLang, ok := loadLangs[lang]; ok == true {
		return loadLang.i18n.Tr(str)
	}
	return ""
}

//实例化
func init() {
	//加载语言包
	langs := dzgin.AppConfig.Strings("lang")
	lang_path := dzgin.AppConfig.String("lang.path")
	if lang_path == "" {
		lang_path = "lang"
	}
	loadLangs = make(map[string]*Load)

	defaultLang = dzgin.AppConfig.String("lang.default")

	basePath := ""
	if utils.FileExists(filepath.Join(dzgin.WorkPath, lang_path)) {
		basePath = dzgin.WorkPath
	} else {
		basePath = dzgin.AppPath
	}
	for _, lang := range langs {
		langfile := filepath.Join(basePath, lang_path, dzgin.AppConfig.String("lang."+lang))
		if err := i18n.SetMessage(lang, langfile); err != nil {
			log.Fatalf("load % Error:%v", langfile, err)
		}

		loadLangs[lang] = &Load{
			&i18n.Locale{Lang: lang},
		}
	}
}
