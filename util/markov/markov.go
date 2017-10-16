package markov

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type Markov struct {
	states map[[2]string][]string
}

func New() *Markov {
	return &Markov{}
}

func (m *Markov) ReadText(text string) string {
	m.Parse(text)

	return m.Generate()
}

func (m *Markov) ReadFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	text, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	m.Parse(string(text))

	return m.Generate()
}

func (m *Markov) ReadURL(URL string) string {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		text := s.Find("p").Text()
		m.Parse(text)
	})

	return m.Generate()
}

func (m *Markov) StateDictionary() map[[2]string][]string {
	return m.states
}

func (m *Markov) Parse(text string) {
	m.states = make(map[[2]string][]string)

	words := strings.Split(text, " ")

	for i := 0; i < len(words)-2; i++ {
		prefix := [2]string{words[i], words[i+1]}

		if _, ok := m.states[prefix]; ok {
			m.states[prefix] = append(m.states[prefix], words[i+2])
		} else {
			m.states[prefix] = []string{words[i+2]}
		}
	}
}

func (m *Markov) Generate() string {
	var sentence bytes.Buffer

	prefix := m.getRandomPrefix([2]string{"", ""})
	sentence.WriteString(strings.Join(prefix[:], " ") + " ")

	for {
		suffix := getRandomWord(m.states[prefix])
		sentence.WriteString(suffix + " ")

		if isTerminalWord(suffix) {
			break
		}

		prefix = [2]string{prefix[1], suffix}
	}

	return sentence.String()
}

func (m *Markov) getRandomPrefix(prefix [2]string) [2]string {
	for key := range m.states {
		if key != prefix {
			prefix = key
			break
		}
	}

	return prefix
}

func getRandomWord(slice []string) string {
	if !(cap(slice) == 0) {
		return slice[rand.Intn(len(slice))]
	} else {
		return ""
	}
}

func isTerminalWord(word string) bool {
	match, _ := regexp.MatchString("(\\.|,|:|;|\\?|!)$", word)
	return match
}
