# Book API

## Table of Contents
* [Running](#running)
* [Structure](#structure)
* [Routes](#routes)
* [Request & Response Examples](#request--response-examples)
* [Testing](#testing)
* [Known Issues](#known-issues)
* [TODO](#todo)

## Running
Docker must be installed on host machine
```bash
cd scripts
sh start.sh - Starts the development docker containers with hot reload
sh start-prod.sh - Starts the production docker containers
sh stop.sh - Tears down the running development containers
sh stop-prod.sh - Tears down the running prod containers
sh tests.sh - Runs the tests and loads the HTML code coverage
sh generate-mocks.sh - Will generate mock classes for unit testing using mockgen
sh docker-purge.sh - Will remove all containers & images
```

## Structure
```
redeam/
 │
 ├──api/                          * Book API source files
 │   ├──domain/                   * Core source files 
 │   └──utils/                    * Helper functions
 │
 ├──mongo-seed/                   * Handles Initializing the local mongodb for local development 
 │   ├──data.json                 * Initial Data
 │   └──Dockerfile                * Creates db and inserts data
 │
 ├──scripts/                      * Folder containing scripts for common tasks  
 ├──vendor/                       * 3rd Party Dependencies 
 ├──.travis.yml                   * Travis CI Configuration
 ├──docker-compose.local.yml      * Docker compose for development
 ├──docker-compose.yml            * Docker compose for production
 ├──Dockerfile                    * Production Dockerfile
 ├──Dockerfile.local              * Development Dockerfile
 ├──Gopkg.toml                    * Describes all the required dependencies
 ├──local.env                     * Environment variables for local development
 └──main.go                       * Program entry point
```

## Routes & Response

| resource                                   | description                       |
|:-------------------------------------------|:----------------------------------|
| [GET /books](#get-books)                   | Returns a list of paginated books, 10 per page |
| [GET /books/[id]](#get-book)               | Returns specified book |
| [POST /books](#post-book)                  | Creates a new book if no duplicate entry based on author, title, publish date fields|
| [PUT /books](#put-book)                    | Updates an existing book |
| [PUT /books/checkout/[id]](#checkout-book) | Checks out a book |
| [PUT /books/checkin/[id]](#checkin-book)   | Checks in a book |
| [PUT /books/[id]/rate/[rate]](#rate-book)  | Rates a book |

## Request & Response Examples
### GET /books
##### Available query params: i.e. books?status=1
* page
* size
* status
* author
* rating

Response body:

    [
      {
        "_id": "5ca7c76f9287bd3832d96f15",
        "author": "Robert Martin",
        "title": "Clean Code",
        "publisher": "Prentice Hall",
        "status": 1,
        "rating": 0,
        "publish_date": "2008"
      },
      {
        "_id": "5ca7c76f9287bd3832d96f16",
        "author": "Paulo Coelho",
        "title": "The Alchemist",
        "publisher": "HarperCollins",
        "status": 1,
        "rating": 2,
        "publish_date": "1988"
      }
    ]

## Testing
Testing uses [GoMock](https://github.com/golang/mock) to generate mocking entities. You will need to download it if updating the test cases.
Once installed, go to domain directory and run mockgen -source=entity_models.go -destination=mock_models.go or run the generate-mocks script in the scripts folder

## TODO

* Update endpoints to return response models
* Add advanced filtering & Sorting using the query models
* Mock database connection
* Finish updating README
