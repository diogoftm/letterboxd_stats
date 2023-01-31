package lbstats

type Gender struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	Name           string `json:"name"`
	Id             int    `json:"id"`
	Logo_path      string `json:"logo_path"`
	Origin_country string `json:"origin_country"`
}

type Country struct {
	Iso_3166_1 string `json:"iso_3166_1"`
	Name       string `json:"name"`
}

type Basic struct {
	Adult                bool      `json:"adult"`
	Backdrop_path        string    `json:"backdrop_path"`
	Budget               int       `json:"budget"`
	Genres               []Gender  `json:"genres"`
	Homepage             string    `json:"homepage"`
	Id                   int       `json:"id"`
	Imdb_id              string    `json:"imdb_id"`
	Original_language    string    `json:"original_language"`
	Original_title       string    `json:"original_title"`
	Overview             string    `json:"overview"`
	Popularity           float32   `json:"popularity"`
	Poster_path          string    `json:"poster_path"`
	Production_companies []Company `json:"production_companies"`
	Production_countries []Country `json:"production_countries"`
	Release_date         string    `json:"release_date"`
	Revenue              int       `json:"revenue"`
	Runtime              int       `json:"runtime"`
	Spoken_languages     []Country `json:"spoken_languages"`
	Status               string    `json:"status"`
	Tagline              string    `json:"tagline"`
	Title                string    `json:"title"`
	Video                bool      `json:"video"`
	Vote_average         float32   `json:"vote_average"`
	Vote_count           int       `json:"vote_count"`
	Year                 int
}

type Cast struct {
	Adult                bool    `json:"adult"`
	Gender               int     `json:"gender"`
	Id                   int     `json:"id"`
	Known_for_department string  `json:"known_for_department"`
	Name                 string  `json:"name"`
	Original_name        string  `json:"original_name"`
	Popularity           float32 `json:"popularity"`
	Profile_path         string  `json:"profile_path"`
	Cast_id              int     `json:"cast_id"`
	Character            string  `json:"character"`
	Credit_id            string  `json:"credict_id"`
	Order                int     `json:"order"`
}

type Crew struct {
	Adult                bool    `json:"adult"`
	Gender               int     `json:"gender"`
	Id                   int     `json:"id"`
	Known_for_department string  `json:"known_for_department"`
	Name                 string  `json:"name"`
	Original_name        string  `json:"original_name"`
	Popularity           float32 `json:"popularity"`
	Profile_path         string  `json:"profile_path"`
	Credit_id            string  `json:"credict_id"`
	Department           string  `json:"department"`
	Job                  string  `json:"job"`
}

type Credits struct {
	Id   int    `json:"id"`
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}

type Rewatch struct {
	Date   string
	Rating float32
}

type Film struct {
	Basic         Basic
	Credits       Credits
	letterboxdURI string
	Date          string
	Rating        float32
	Rewatch       []Rewatch
}

type FilmList struct {
	AllFilms        map[int]*Film    //keys are TMDB id's
	AllFilmsByLbURI map[string]*Film //keys are letterboxd URI's
	FilmsByYear     map[int][]*Film
}

type BasicStats struct {
	NFilms          int
	NRewatched      int
	MostSeenDecades map[int]int
	NMinutes        int
	NWeek           [52]int
	NDayOfWeek      [7]int //0 for sunday, ...
	NMonth          [12]int
	Genres          map[string]int
	Languages       map[string]int
	Countries       map[string]int
	AvgRating       float32
}

type CreditsStats struct {
	Acting           map[string]int
	Directors        map[string]int
	Writers          map[string]int
	Cinematographers map[string]int
	Editors          map[string]int
	Producers        map[string]int
	Music            map[string]int
}

type User struct {
	DateJoined    string
	Username      string
	GivenName     string
	FamilyName    string
	Email         string
	Location      string
	Website       string
	Bio           string
	Pronoun       string
	FavoriteFilms []string //list of Letterboxd URI's (to get film info access AllFilmsLbURI map)
}

// Returns a map with jobs as keys and a list of names as value
func (credits Credits) Jobs() map[string][]string {
	jobs := make(map[string][]string)
	for _, v := range credits.Crew {
		jobs[v.Job] = append(jobs[v.Job], v.Name)
	}
	return jobs
}
