# Movie App
## _A Simple Movie App_

## Features

- Admin user can create new movies
- Guest users can search for movie info, ratings and comments
- Loged-in users can get movies, which already rated or commented
- Logged-in users can rate a movie
- Logged-in users can add comments on a movie

## API Endpoints

- GET -  http://localhost:3000/api/movie/search/{movieName} [public API - Movie search with name. Will return movie info, avg rating, number of ratings and comments]
- GET - http://localhost:3000/api/movie/list [private API - Movie search with userId and movie Id. Will return movie info, user's rating and comments]
- POST - http://localhost:3000/api/movie [private API (Admin only) - Create new Movie into the database. Will return the create result]
- PUT - http://localhost:3000/api/movie/rating/{id} [private API - Rate a movie with userId and movieId (movie Id can be obtained from public search). Will return the status of rating]
- PUT - http://localhost:3000/api/movie/comment/{id} [private API - Add comment to a movie with userId and movieId (movie Id can be obtained from public search). Will return the status of addition of new comment]
- 

POSTMAN collection and db dump are included in the repo

mongodb import dump : mongorestore -d movies_db <path>
Also execute the index : db.getCollection('movie').createIndex({ name: "text" })
