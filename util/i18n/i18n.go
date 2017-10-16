package i18n
	
//import "fmt"
import 	"github.com/nicksnyder/go-i18n/i18n"
import 	"github.com/NamedKitten/KittehBotGo/util/static"

func init() {
	file, err := static.ReadFile("translations/en-gb.yaml")
	if (err != nil) {
		panic(err)
	}
	i18n.ParseTranslationFileBytes("en-gb.yaml", file)
	file, err = static.ReadFile("translations/fr-fr.yaml")
	if (err != nil) {
		panic(err)
	}
	i18n.ParseTranslationFileBytes("fr-fr.yaml", file)
	file, err = static.ReadFile("translations/es-es.yaml")
	if (err != nil) {
		panic(err)
	}
	i18n.ParseTranslationFileBytes("es-es.yaml", file)
	//fmt.Println("Initialising map.")
	//Translations = make(map[string]map[string]string) 
}
