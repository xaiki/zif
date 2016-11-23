package data

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sajari/fuzzy"
	log "github.com/sirupsen/logrus"
)

// This provides searching, as it is a little more comlex than just a db query.
// Search strings could provide other data that needs parsing, as well as spell
// correction that needs doing. This has to be passed through other functions
// before it hits a db query, hence this.
type SearchProvider struct {
	Loaded bool

	model *fuzzy.Model
	// if the model has been loaded, otherwise no autocomplete/spell suggestions
}

func NewSearchProvider() *SearchProvider {
	sp := &SearchProvider{false, nil}

	sp.LoadModel()

	return sp
}

func (sp *SearchProvider) SaveModel() error {
	if !sp.Loaded {
		return errors.New("Model has not been loaded, save failed")
	}

	return sp.model.Save("./data/model.dat")
}

// Loads the model from disk, if it does not exist then load the raw corpus.
func (sp *SearchProvider) LoadModel() {
	var err error

	// Train with a corpus if the model has not already been built and saved.
	// Popular torrents will also be added to this.
	if _, err = os.Stat("./data/model.dat"); os.IsNotExist(err) {
		err = sp.loadCorpus()

		if err != nil {
			log.Error(err.Error())
		}

		return
	}

	sp.model, err = fuzzy.Load("./data/model.dat")

	if err != nil {
		return
	}

	sp.Loaded = true
}

func (sp *SearchProvider) loadCorpus() error {
	log.Info("Model does not exist, loading corpus.")

	if _, err := os.Stat("./data/corpus.txt"); os.IsNotExist(err) {
		return err
	}

	corpus, err := os.Open("./data/corpus.txt")

	if err != nil {
		return err
	}

	// loop through all words, train the model by these.
	scanner := bufio.NewScanner(corpus)
	scanner.Split(bufio.ScanWords)

	sp.model = fuzzy.NewModel()

	for scanner.Scan() {
		sp.model.TrainWord(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	sp.Loaded = true

	return sp.SaveModel()
}

func (sp *SearchProvider) spellCheck(query string) string {
	scanner := bufio.NewScanner(strings.NewReader(query))
	scanner.Split(bufio.ScanWords)

	newQuery := bytes.Buffer{}

	for scanner.Scan() {
		newQuery.WriteString(sp.model.SpellCheck(scanner.Text()))
		newQuery.WriteString(" ")
	}

	// Remove the space at the end
	newQuery.Truncate(newQuery.Len() - 1)

	return newQuery.String()
}

func IsAlnumWord(word string) bool {
	for _, i := range word {
		if !unicode.IsLetter(i) && !unicode.IsNumber(i) {
			return false
		}
	}

	return true
}

// Takes a string, makes it look "nice" for an autocomplete cue.
func SanitiseForAuto(in string) string {
	buffer := bytes.Buffer{}

	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if IsAlnumWord(scanner.Text()) {
			buffer.WriteString(scanner.Text())
			buffer.WriteString(" ")
		}
	}

	return buffer.String()
}

func (sp *SearchProvider) Suggest(db *Database, query string) ([]string, error) {
	checked, err := db.Suggest(fmt.Sprintf("%s%%", sp.spellCheck(query)))

	if err != nil {
		return nil, err
	}

	nonChecked, err := db.Suggest(fmt.Sprintf("%s%%", query))
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(checked)+len(nonChecked))

	for _, i := range checked {
		ret = append(ret, SanitiseForAuto(i))
	}

	for _, i := range nonChecked {
		ret = append(ret, SanitiseForAuto(i))
	}

	return ret, nil
}

func (sp *SearchProvider) Search(db *Database, query string, page int) ([]*Post, error) {
	// TODO: Instead of searching for spell-corrected versions, suggest an
	// alternate search.
	results, err := db.Search(query, page, 25)

	return results, err
}
