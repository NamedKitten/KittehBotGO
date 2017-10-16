package i18n
	
//import "fmt"
import 	"github.com/nicksnyder/go-i18n/i18n"
import 	"github.com/NamedKitten/KittehBotGo/util/static"

func init() {
	files, _ := static.WalkDirs("translations/", false)
	for _, filename := range files {
		file, err := static.ReadFile(filename)
		if (err != nil) {
			panic(err)
		}
		i18n.ParseTranslationFileBytes(filename, file)
		
	}
	//fmt.Println("Initialising map.")
	//Translations = make(map[string]map[string]string) 
}
