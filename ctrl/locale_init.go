package ctrl

import (
	"os"
	"path"

	"github.com/liangx8/i18n"
)

var langs = [2]string{"zh_CN", "en_US"}
var langDesc = [2]string{"Simplified Chinese", "Englis(US)"}

// @prefix root directory of internationalization
func LocaleInit(prefix string) error {
	for ix, lang := range langs {
		if err := i18n.Register(lang, langDesc[ix], i18n.NewJsonDecoder()); err != nil {
			return err
		}
		langPath := path.Join(prefix, lang)
		entrs, err := os.ReadDir(langPath)
		if err != nil {
			return err
		}
		for _, ent := range entrs {
			if ent.Type().IsRegular() {
				if err := i18n.AddResource(lang, i18n.ResourceFilename(path.Join(langPath, ent.Name()))); err != nil {
					return err
				}
			}

		}
	}
	i18n.SetLang(langs[0])
	return nil
}
