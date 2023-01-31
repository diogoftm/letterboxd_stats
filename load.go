package lbstats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type FindAnswer struct {
	Results       []Basic `json:"results"`
	Total_results int     `json:"total_results"`
}

// List all films on a FilmList
func (lst FilmList) ListFilms() []*Film {
	var filmList []*Film
	for _, v := range lst.AllFilms {
		filmList = append(filmList, v)
	}
	return filmList
}

// Load films from diary.csv like file
func LoadFilmsFromCSVfiles(diaryPath string) FilmList {
	channel := make(chan []interface{})
	wg := new(sync.WaitGroup)

	in1, err := os.ReadFile(diaryPath)
	Check(err)

	r := csv.NewReader(strings.NewReader(string(in1)))
	r.Read()

	replacer := strings.NewReplacer(" ", "+", "#", "%23", "&", "%26")
	var rating float32
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		Check(err)
		if record[4] == "" {
			rating = 0
		} else {
			rating = SingleParseFloat(record[4])
		}
		rewatch := false
		if record[5] == "Yes" {
			rewatch = true
		}
		wg.Add(1)
		go func(title string, year string, lbURI string, date string, rewatch bool, rating float32, wgp *sync.WaitGroup) {
			defer wgp.Done()
			f := GetFilm(replacer.Replace(title), year)
			if f == nil {
				return
			}
			f.Date = date
			f.Basic.Year = SingleAtoi(record[2])
			f.Rating = rating
			channel <- []interface{}{f.Basic.Id, SingleAtoi(date[0:4]), f, rewatch, lbURI}
		}(record[1], record[2], record[3], record[7], rewatch, rating, wg)
	}

	allFilms := make(map[int]*Film)
	allFilmsLB := make(map[string]*Film)
	filmsByYear := make(map[int][]*Film)
	wgc := new(sync.WaitGroup)
	wgc.Add(1)
	go func() {
		defer wgc.Done()
		for film := range channel {
			v, found := allFilms[film[0].(int)]
			if !found {
				if film[3].(bool) {
					film[2].(*Film).Rewatch = append(film[2].(*Film).Rewatch, Rewatch{film[2].(*Film).Date, film[2].(*Film).Rating})
				}
				allFilms[film[0].(int)] = film[2].(*Film)
				allFilmsLB[film[4].(string)] = film[2].(*Film)
				filmsByYear[film[1].(int)] = append(filmsByYear[film[1].(int)], film[2].(*Film))
			} else {
				if film[3].(bool) {
					v.Rewatch = append(v.Rewatch, Rewatch{film[2].(*Film).Date, film[2].(*Film).Rating})
					filmsByYear[film[1].(int)] = append(filmsByYear[film[1].(int)], v)
				} else {
					v.Date = film[2].(*Film).Date
					v.Rating = film[2].(*Film).Rating
				}
			}
		}
	}()

	wg.Wait()
	close(channel)
	wgc.Wait()
	return FilmList{allFilms, allFilmsLB, filmsByYear}
}

// Load user info from profile.csv like file
func LoadUser(profilePath string) User {
	var user User
	in1, err := os.ReadFile(profilePath)
	Check(err)

	r := csv.NewReader(strings.NewReader(string(in1)))
	r.Read()

	record, err := r.Read()
	Check(err)
	user.DateJoined = record[0]
	user.Username = record[1]
	user.GivenName = record[2]
	user.FamilyName = record[3]
	user.Email = record[4]
	user.Location = record[5]
	user.Website = record[6]
	user.Bio = record[7]
	user.Pronoun = record[8]

	for it := 9; it < len(record); it++ {
		user.FavoriteFilms = append(user.FavoriteFilms, record[it])
	}

	return user
}

// TMDB API requests for film info
func GetFilm(title string, year string) *Film {
	apiKey := "c54aa98d0ceea0ec6aba9ab2ece643ee"
	url1 := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&year=%s", apiKey, title, year)

	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic("Error build request")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic("Error while connecting with TMDB")
	}

	defer resp.Body.Close()

	var f Film
	var b FindAnswer
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		log.Println(err)
	}

	if len(b.Results) == 0 {
		return nil
	}
	id := b.Results[0].Id

	url2 := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s", id, apiKey)

	req, err = http.NewRequest("GET", url2, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic("Error build request")
	}

	client = &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic("Error while connecting with TMDB")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&f.Credits); err != nil {
		log.Println(err)
	}

	url3 := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", id, apiKey)

	req, err = http.NewRequest("GET", url3, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		panic("Error build request")
	}

	client = &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		panic("Error while connecting with TMDB")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&f.Basic); err != nil {
		log.Println(err)
	}

	return &f
}
