package structs

import "net/url"

type Movie struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Year     int    `json:"year" bson:"year"`
	Director string `json:"director" bson:"director"`
	Genre    string `json:"genre" bson:"genre"`
}

type MovieSearchResult struct {
	ID            string   `json:"id" bson:"_id"`
	Name          string   `json:"name" bson:"name"`
	Year          int      `json:"year" bson:"year"`
	Director      string   `json:"director" bson:"director"`
	Genre         string   `json:"genre" bson:"genre"`
	Rating        float32  `json:"rating,omitempty" bson:"rating,omitempty"`
	AverageRating float32  `json:"averageRating,omitempty" bson:"averageRating,omitempty"`
	TotalRatings  int64    `json:"totalRating,omitempty" bson:"totalRating,omitempty"`
	Comments      []string `json:"comments" bson:"comments"`
}

type Rating struct {
	ID      string  `json:"id,omitempty" bson:"_id,omitempty"`
	MovieID string  `json:"movieId,omitempty" bson:"movieId"`
	Rating  float32 `json:"rating" bson:"rating"`
	UserId  string  `json:"userId,omitempty" bson:"userId"`
}

type Comment struct {
	ID      string `json:"_id,omitempty" bson:"_id,omitempty"`
	MovieID string `json:"movieId,omitempty" bson:"movieId"`
	Comment string `json:"comment" bson:"comment"`
	UserId  string `json:"userId,omitempty" bson:"userId"`
}

type User struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
}

type AVGRating struct {
	ID      string  `json:"id" bson:"_id"`
	Average float32 `json:"average" bson:"average"`
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func (movie *Movie) Validate() url.Values {
	errs := url.Values{}

	// check if the ID empty
	if movie.ID == "" {
		errs.Add("id", "The id field is required!")
	}

	if movie.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	if movie.Year == 0 {
		errs.Add("year", "The year field is required!")
	}

	if movie.Genre == "" {
		errs.Add("genre", "The genre field is required!")
	}

	return errs
}
