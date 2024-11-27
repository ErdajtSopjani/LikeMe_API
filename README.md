<div align="center">
  <img src="https://github.com/user-attachments/assets/fd4e29ff-aca9-4a37-ad11-d09a4097e9af" alt="likeme-icon" width="200" />
</div>

# LikeMe_API

Back-end of the LikeMe social network platform

Brief description of LikeMe:

***LikeMe is a social media network made to connect Like-minded people. My goal with LikeMe is to make an AI powered social media application in the simplest way possible, such that even a beginner can understand the code. The code will be open-source and available on my GitHub. Although the project is not yet finished, I'm working on it on my free time.***

## Technical Description

LikeMe_API is the simplest possible backend for a social media network implemented in go, it uses `Postgres` as the database. It's designed to be a REST-API and I've implemented unit-tests for each endpoint that check numerous request scenarios and if the server returns the expected body & code.

### Technologies:

*Language: `Go`*

*Database: `Postgres`*

*ORM: `gorm`*

*Routes, requests etc: `go-chi`*

*Auth: `self-made email OTP`*

*Email: `sendgrid`*

Everything else is implemented by the go standard library.

### Authentication:

I've made a simple OTP email authentication system, since It's the simplest one to implement by myself and the cheapest overall.
It works by sending email 2fa codes and when used it will provide a token assigned to a specific user which will be used to authenticate requests coming to the API.

### Goal:

My goal was to make the simplest possible backend for a social media platform, the code is simple and focused on readability.

### Running:

1. Install all the packages with go
2. `go run *.go` will run the server

The server will run on `:8080`

### Tests:

I've implemented unit tests for each enpoint, this will test multiple variations on the request and will check if the excpected body & code match.
