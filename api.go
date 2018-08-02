package urbandictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SearchResult contains data regarding an Urban Dictionary search result
type SearchResult struct {
	Results []Result `json:"list"`
}

// Result contains data regarding an Urban Dictionary definition
type Result struct {
	Author       string    `json:"author"`
	CurrentVote  string    `json:"current_vote"`
	Date         time.Time `json:"-"`
	Definition   string    `json:"definition"`
	DefinitionID int       `json:"defid"`
	Example      string    `json:"example"`
	Permalink    string    `json:"permalink"`
	ThumbsUp     int       `json:"thumbs_up"`
	ThumbsDown   int       `json:"thumbs_down"`
	Word         string    `json:"word"`

	WrittenOn string `json:"written_on"`
}

// UrbanDictionaryAPI contains the Urban Dictionary API URL for search results
const UrbanDictionaryAPI = "http://api.urbandictionary.com/v0/define?term="

// Query returns data from an Urban Dictionary search result
func Query(searchTerm string) (*SearchResult, error) {
	resp, err := http.Get(UrbanDictionaryAPI + url.QueryEscape(searchTerm))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Response was not a 200: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &SearchResult{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, errors.New("No results were found")
	}

	for _, result := range res.Results {
		time, err := time.Parse("2006-01-02T15:04:05.000Z", result.WrittenOn)
		if err != nil {
			return nil, err
		}
		result.Date = time
	}

	return res, nil
}
