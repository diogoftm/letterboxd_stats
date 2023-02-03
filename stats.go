package lbstats

import (
	"errors"
	"time"
)

// General stats based on a list of films for a given year (year=0 for no year restriction)
func GetBasicStats(list FilmList, year int) (BasicStats, error) {
	var b BasicStats
	var films []*Film

	if year == 0 {
		films = list.ListFilms()
		b.NFilms = len(list.AllFilms)
		fr, _, counter := nRewatched(list, year)
		b.NRewatched = counter
		films = append(films, fr...) //add rewatched films
	} else {
		f, ok := list.FilmsByYear[year]
		if ok {
			films = f
			fr, countUnique, counter := nRewatched(list, year)
			b.NRewatched = counter
			films = append(films, Duplicates(fr)...) //add films that were rewatched more than once
			b.NFilms = len(f) - countUnique
		} else {
			return b, errors.New("Invalid year")
		}
	}

	b.MostSeenYears = make(map[int]int)
	b.Genres = make(map[string]int)
	b.Languages = make(map[string]int)
	b.Countries = make(map[string]int)

	var dateWatched time.Time
	nRated := 0
	for _, v := range films {
		dateWatched = ParseToDate(v.Date)
		b.NMinutes += v.Basic.Runtime
		b.NMonth[dateWatched.Month()-1]++
		b.NDayOfWeek[dateWatched.Weekday()]++
		_, nWeek := dateWatched.ISOWeek()
		b.NWeek[nWeek-1]++

		if v.Rating != 0 {
			b.AvgRating += v.Rating
			nRated++
		}
		for _, g := range v.Basic.Genres {
			b.Genres[g.Name]++
		}

		b.Languages[v.Basic.Original_language]++

		for _, c := range v.Basic.Production_countries {
			b.Countries[c.Name]++
		}
		b.MostSeenYears[v.Basic.Year]++
	}
	b.AvgRating /= float32(nRated)
	return b, nil
}

// Cast and crew stats based on a list of films for a given year (year=0 for no year restriction)
func GetCreditsStats(list FilmList, year int) (CreditsStats, error) {
	var c CreditsStats
	var films []*Film

	if year == 0 {
		films = list.ListFilms()
	} else {
		f, ok := list.FilmsByYear[year]
		if ok {
			films = f
		} else {
			return c, errors.New("Invalid year")
		}
	}

	c.Acting = make(map[string]int)
	c.Cinematographers = make(map[string]int)
	c.Directors = make(map[string]int)
	c.Editors = make(map[string]int)
	c.Writers = make(map[string]int)
	c.Producers = make(map[string]int)
	c.Music = make(map[string]int)

	for _, f := range films {
		for _, castMember := range f.Credits.Cast {
			c.Acting[castMember.Name]++
		}
		for job, names := range f.Credits.Jobs() {
			for _, name := range names {
				switch job {
				case "Director":
					c.Directors[name]++
				case "Editor":
					c.Editors[name]++
				case "Writer":
					c.Writers[name]++
				case "Screenplay":
					c.Writers[name]++
				case "Cinematography":
					c.Cinematographers[name]++
				case "Producer":
					c.Producers[name]++
				case "Music":
					c.Music[name]++
				}
			}
		}
	}

	return c, nil
}

// number of rewatches
func nRewatched(list FilmList, year int) ([]*Film, int, int) {
	counter := 0
	counterUnique := 0
	var v []*Film
	var rewatched []*Film
	if year == 0 {
		v = list.ListFilms()
		for _, f := range v {
			counter += len(f.Rewatch)
			for it := 0; it < len(f.Rewatch); it++ {
				rewatched = append(rewatched, f)
				if it == 0 {
					counterUnique++
				}
			}
		}
	} else {
		v = list.FilmsByYear[year]
		flag := 0
		for _, f := range v {
			flag = 0
			for _, r := range f.Rewatch {
				if SingleAtoi(r.Date[0:4]) == year {
					counter++
					rewatched = append(rewatched, f)
					if flag == 0 {
						counterUnique++
					}
					flag = 1
				}
			}
		}
	}
	return rewatched, counterUnique, counter
}
