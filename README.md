# Project Title

Golang Web Example

## Description

A basic Rest API With JWT authorization and Redis caching.



## Features And Info

* Postgres Database
* JWT Authorization
* Chi Router
* Redis Caching
* SQLC Queries


## Routes

- /auth/*** - Handling Authentication
- /user/*** - Handling User Services


## /auth
- /signin - Sign In With Existing User
- /signup - Register User
- /refresh - Refresh JWT Token

## /user

- /all - Get All Users
- /{id} - Get User By Id
- /delete/{id} - Delete User By Id
- /home - Home Page (Get User Credentials Using JWT Token)