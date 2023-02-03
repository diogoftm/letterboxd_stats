# letterboxd_stats

Package to load films, user information and create stats based on [Letterboxd](https://letterboxd.com/) .csv exportable files.

## Stats
> When calling `LoadFilmsFromCSVfile()` a list of films is returned, using it anything can be done with the film info privided
by TMDB. So, if you wish to create your own stats or just display the list of films there is no need to get the buil-in stats. 
This package can be used just for loading data from the .csv files and obtain the information about each film.

### Basic Stats
- Number of films - `NFilms int`
- Number of rewatched films - `NRewatched int`
- Most seen years - `MostSeenYears map[int]int`
- Total of minutes - `NMinutes int`
- Number of films watched per week - `NWeek [52]int`
- Number of films watched per fay of week (0 for sunday, ...) - `NDayOfWeek [7]int`
- Number of films watched per month - `NMonth [12]int`
- Number of films per genre - `Genres map[string]int`
- Number of films per languages - `Languages map[string]int`
- Number of films per country - `Countries map[string]int`
- Average rating - `AvgRating float32`

### Cast Stats
- Number of films per cast member - `Acting map[string]int`
- Number of films per director - `Directors map[string]int`
- Number of films per writer - `Writers map[string]int`
- Number of films per cinematographer - `Cinematographers map[string]int`
- Number of films per editor - `Editors map[string]int`
- Number of films per producer - `Producers map[string]int`
- Number of films per music crew member - `Music map[string]int`

## Example
```go
import "github.com/diogoftm/letterboxd_stats"
// get user info
var user lbstats.User
user = lbstats.LoadUser("xpto/profile.csv") // get user info

// load films
var filmList lbstats.FilmList
filmList = lbstats.LoadFilmsFromCSVfile("xpto/diary.csv")

// get basic and credits stats
var bs lbstats.BasicStats
var cs lbstats.CreditsStats
var err error
bs, err = lbstats.GetBasicStats(filmList, 0)
cs, err = lbstats.GetCreditsStats(filmList, 0)
```

## TMDB
>[The Movie Database (TMDB)](https://www.themoviedb.org/) is a community built movie and TV database. Every piece of data has been added by our amazing community dating back to 2008. TMDB's strong international focus and breadth of data is largely unmatched and something we're incredibly proud of. Put simply, we live and breathe community and that's precisely what makes us different.

All film information is obtained from TMDB's API. Letterboxd uses TMDB too, but on the .csv files they do not include the TMDB id of films. So, a search is always made using the title and the year of release of each film and the best match is picked. Because of that, in some cases, a Letterboxd URI might be mismatched to it's respective TMDB id.
