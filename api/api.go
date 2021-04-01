package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"movies-app/helpers"
	"movies-app/structs"
	"net/http"
	"strings"
)

//public API - Movie search with name.
//Will return movie info, avg rating, number of ratings and comments
func PublicSearch(response http.ResponseWriter, request *http.Request) {
	movieName, ok := mux.Vars(request)["name"]
	//movieId := params[0]
	if !ok || len(movieName) < 1 {
		err := map[string]interface{}{"validationError": "please provide valid movie name search key"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	fmt.Println("Movie search with guest/public user")
	var movies []structs.Movie
	var movieIds []string
	// connect db
	collection := helpers.GlobalMongoClient.Collection("movie")
	query := bson.M{"$text": bson.M{"$search": movieName}}
	cur, err := collection.Find(context.TODO(), query)
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var movie structs.Movie
		err := cur.Decode(&movie)
		if err != nil {
			fmt.Println(err)
		}
		movieIds = append(movieIds, movie.ID)
		movies = append(movies, movie)
	}

	var avgRating []structs.AVGRating
	matchStage := bson.D{
		{Key: "$match", Value: bson.M{"movieId": bson.M{"$in": movieIds}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$movieId"},
			{Key: "average", Value: bson.D{
				{Key: "$avg", Value: "$rating"}}}}}}
	pipeline := mongo.Pipeline{matchStage, groupStage}
	collection = helpers.GlobalMongoClient.Collection("rating")
	cur, err = collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	if err = cur.All(context.TODO(), &avgRating); err != nil {
		helpers.GetError(err, response)
		return
	}
	var comments []structs.Comment
	collection = helpers.GlobalMongoClient.Collection("comments")
	cur, err = collection.Find(context.TODO(), bson.M{"movieId": bson.M{"$in": movieIds}})
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var comment structs.Comment
		err := cur.Decode(&comment)
		if err != nil {
			fmt.Println(err)
		}
		comments = append(comments, comment)
	}
	var movieSearchList []structs.MovieSearchResult
	collection = helpers.GlobalMongoClient.Collection("rating")
	for _, mId := range movieIds {
		var movieSearch structs.MovieSearchResult
		movieSearch.ID = mId
		for _, movie := range movies {
			if movie.ID == mId {
				movieSearch.Name = movie.Name
				movieSearch.Year = movie.Year
				movieSearch.Director = movie.Director
				movieSearch.Genre = movie.Genre
				break
			}
		}
		var commentList []string
		for _, comment := range comments {
			if comment.MovieID == mId {
				commentList = append(commentList, comment.Comment)
			}
		}
		movieSearch.Comments = commentList
		for _, rating := range avgRating {
			if rating.ID == mId {
				movieSearch.AverageRating = rating.Average
				break
			}
		}
		query := bson.M{"movieId": mId}
		count, err := collection.CountDocuments(context.TODO(), query)
		if err != nil {
			helpers.GetError(err, response)
			return
		}
		movieSearch.TotalRatings = count
		movieSearchList = append(movieSearchList, movieSearch)
	}
	fmt.Println("Movie search with guest/public user completed")
	helpers.OK(movieSearchList, response)
}

//private API - Movie search with userId and movie Id.
//Will return movie info, user's rating and comments
func SearchByUserId(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	fmt.Println("Movie search with user logged in : ", user)
	var ratings []structs.Rating
	var movieIds []string
	// connect db
	collection := helpers.GlobalMongoClient.Collection("rating")
	cur, err := collection.Find(context.TODO(), bson.M{"userId": user})
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var rating structs.Rating
		err := cur.Decode(&rating)
		if err != nil {
			fmt.Println(err)
		}
		movieIds = append(movieIds, rating.MovieID)
		ratings = append(ratings, rating)
	}
	var movies []structs.Movie
	collection = helpers.GlobalMongoClient.Collection("movie")
	cur, err = collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": movieIds}})
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var movie structs.Movie
		err := cur.Decode(&movie)
		if err != nil {
			fmt.Println(err)
		}
		movies = append(movies, movie)
	}
	var comments []structs.Comment
	collection = helpers.GlobalMongoClient.Collection("comments")
	cur, err = collection.Find(context.TODO(), bson.M{"movieId": bson.M{"$in": movieIds}, "userId": user})
	if err != nil {
		helpers.GetError(err, response)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var comment structs.Comment
		err := cur.Decode(&comment)
		if err != nil {
			fmt.Println(err)
		}
		comments = append(comments, comment)
	}
	var movieSearchList []structs.MovieSearchResult
	for _, mId := range movieIds {
		var movieSearch structs.MovieSearchResult
		movieSearch.ID = mId
		for _, movie := range movies {
			if movie.ID == mId {
				movieSearch.Name = movie.Name
				movieSearch.Year = movie.Year
				movieSearch.Director = movie.Director
				movieSearch.Genre = movie.Genre
				break
			}
		}
		var commentList []string
		for _, comment := range comments {
			if comment.MovieID == mId {
				commentList = append(commentList, comment.Comment)
			}
		}
		movieSearch.Comments = commentList
		for _, rating := range ratings {
			if rating.MovieID == mId {
				movieSearch.Rating = rating.Rating
				break
			}
		}
		movieSearchList = append(movieSearchList, movieSearch)
	}
	fmt.Println("Movie search with logged in user completed")
	helpers.OK(movieSearchList, response)
}

//private API (Admin only) - Create new Movie into the database.
//Will return the create result
func CreateMovie(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	fmt.Println("Create Movie request recieved from user : ", user)
	response.Header().Set("Content-Type", "application/json")
	var newMovie structs.Movie
	// decode our body request params
	_ = json.NewDecoder(request.Body).Decode(&newMovie)
	//Validate payload
	if validErrs := newMovie.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	// connect db
	collection := helpers.GlobalMongoClient.Collection("movie")
	// insert movie model
	_, err := collection.InsertOne(context.TODO(), newMovie)
	if err != nil {
		fmt.Println("Created new Movie request failed with error : ", err.Error())
		if strings.Contains(err.Error(), "E11000") {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(400)
			response.Write([]byte(`{"status":"Duplicate movie Id"}`))
			return
		}
		helpers.GetError(err, response)
		return
	}
	fmt.Println("Created new Movie with id : ", newMovie.ID)
	helpers.OK(newMovie, response)
}

//private API - Rate a movie with userId and movieId (movie Id can be obtained from public search).
//Will return the status of rating
func RateMovie(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	movieId, ok := mux.Vars(request)["id"]
	//movieId := params[0]
	if !ok || len(movieId) < 1 {
		err := map[string]interface{}{"validationError": "please provide valid movieId"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	fmt.Println("Rate Movie request recieved from user : ", user, "for movieId: ", movieId)
	response.Header().Set("Content-Type", "application/json")
	client := helpers.GlobalMongoClient
	query := bson.M{
		"_id": movieId}
	c := client.Collection("movie")
	ctx := context.TODO()
	var movie structs.Movie
	c.FindOne(ctx, query).Decode(&movie)
	if movie.ID == "" {
		err := map[string]interface{}{"validationError": "please provide valid movieId"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	var rating structs.Rating
	// decode our body request params
	_ = json.NewDecoder(request.Body).Decode(&rating)
	//Validate payload
	if rating.Rating > 10 || rating.Rating < 1 {
		err := map[string]interface{}{"validationError": "rating range is 1-10"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	// connect db
	collection := helpers.GlobalMongoClient.Collection("rating")
	// update rating model
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"movieId": movieId}
	update := bson.M{
		"$set": bson.M{
			"rating": rating.Rating,
			"userId": user,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		fmt.Println("rate a movie request failed with error : ", err.Error())
		helpers.GetError(err, response)
		return
	}
	fmt.Println("rated the movie : ", movieId, " with rating  : ", rating.Rating, " by user : ", user)
	helpers.OK("", response)
}

//private API - Add comment to a movie with userId and movieId (movie Id can be obtained from public search).
//Will return the status of addition of new comment
func CommentMovie(response http.ResponseWriter, request *http.Request) {
	user, _, _ := request.BasicAuth()
	movieId, ok := mux.Vars(request)["id"]
	//movieId := params[0]
	if !ok || len(movieId) < 1 {
		err := map[string]interface{}{"validationError": "please provide valid movieId"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	fmt.Println("Comment Movie request recieved from user : ", user, "for movieId: ", movieId)
	response.Header().Set("Content-Type", "application/json")
	if movieId == "" {
		err := map[string]interface{}{"validationError": "please provide valid movieId"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	client := helpers.GlobalMongoClient
	query := bson.M{
		"_id": movieId}
	c := client.Collection("movie")
	ctx := context.TODO()
	var movie structs.Movie
	c.FindOne(ctx, query).Decode(&movie)
	if movie.ID == "" {
		err := map[string]interface{}{"validationError": "please provide valid movieId"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	var comment structs.Comment
	// decode our body request params
	_ = json.NewDecoder(request.Body).Decode(&comment)
	//Validate payload
	if comment.Comment == "" {
		err := map[string]interface{}{"validationError": "please write something"}
		response.WriteHeader(400)
		helpers.OK(err, response)
		return
	}
	// connect db
	collection := helpers.GlobalMongoClient.Collection("comments")
	// update comment model
	comment.MovieID = movieId
	comment.UserId = user
	_, err := collection.InsertOne(context.TODO(), comment)
	if err != nil {
		fmt.Println("Adding new comment for a Movie request failed with error : ", err.Error())
		helpers.GetError(err, response)
		return
	}
	fmt.Println("Comment added for Movie with id : ", movieId, "by user : ", user)
	helpers.OK("", response)
}
