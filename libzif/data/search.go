package data

import (
	"bufio"
	"errors"
	"os"

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

func (sp *SearchProvider) Suggest(query string) {
	// TODO: Implement this...
	/*
		Add popular post titles to the corpus, then see if I can have the model
		suggest entire titles.
		It is unlikely though...
		The best thing to do is probably spell correct the search, then perform
		a tiny search on the database and return the results live.
	*/
}

func (sp *SearchProvider) Search(db *Database, query string, page int) ([]*Post, error) {
	// TODO: Spell correct query, then search with the corrected version.
	results, err := db.Search(query, page)

	return results, err
}
