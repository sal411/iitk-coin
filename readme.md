# IITK COINS

## Structure
### Controllers
- checkAPI
    * Sends a message that API is LIVE 
    * GET method
- createuser
    * Creates a User in database, with provied name, roll number and password
    * POST method
- login
    * Logs in the user, checks is user exists in database
    * creates a JWT Token, and sets it in cookie
    * POST method
- secretpage
    * checks if users is logged in, then gives access to data
    * GET Method
- user
    * Has functions pertaining to user, create database, add user to database
- bank
    * Has function to handle getCoin, transferCoins, addCoins requests
    

### models
- token
- userdata
- bankdata
- transferCoin

### database
- bankDataBase
    * creates a bank data base, and opens account for a user
- userDataBase
    * creates a user data base, and adds user
- getFromBankDB
    * has functions for reading data from bank database
- getFromUserDB
    * has functions for reading data from user database
- postToBankDV
    * has functions for wriring data from user database

### routes
- routes
    - handles all the http requests, and routes them to appropriate functions

### utils
- db
    - gets a connection to database
- generateToken
    - generates a JWT token for a given user
- getFromDataBase
    - has all functions to fetch data from database
- printError
- verifyToken
    - verifies the JWT Token

### main.go

## To Run this app, perform the following step in order

1. Clone this repo to your machine
2. cd into the project folder
3. enter `go run main.go`  to start server



 