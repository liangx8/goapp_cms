package ctrl

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/unknwon/i18n"
)

var langs = [2]string{"zh_CN", "en_US"}
var langDesc = [2]string{"Chinese", "Englis(US)"}

func LocaleInit(prefix string) error {
	files, err := os.ReadDir(prefix)
	if err != nil {
		return err
	}
	for ix, lang := range langs {
		msgs := make([]any, 0)
		for _, file := range files {
			if strings.Contains(file.Name(), lang) {
				msgs = append(msgs, path.Join(prefix, file.Name()))
			}
		}
		if len(msgs) == 1 {
			if err := i18n.SetMessageWithDesc(lang, langDesc[ix], msgs[0]); err != nil {
				log.Print(err)
			}
		}
		if len(msgs) > 1 {
			if err := i18n.SetMessageWithDesc(lang, langDesc[ix], msgs[0], msgs[1:]...); err != nil {
				log.Print(err)
			}
		}
	}

	//	i18n.SetMessageWithDesc("zh-CN", "Chinese", "messages/locale_zh_CN.ini")
	//	i18n.SetMessageWithDesc("en-US", "English(US)", "messages/locale_en_US.ini")

	log.Print(i18n.Tr("zh_CN", "hi", "旅行者"))
	log.Print(i18n.Tr("en_US", "hi", "旅行者"), i18n.Tr("zh_CN", "login"))
	return nil
}
