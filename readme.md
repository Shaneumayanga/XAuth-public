# How it works

```bash
git clone https://github.com/Shaneumayanga/XAuth-public
```

## start oAuth flow :-
http://localhost:8080/login/oauth/authorize?client_id=84bd6d67-892d-4ebb-ab8f-361504388042&redirect_url=https://shaneumayanga.com/redirect&response_type=code
 the above url with give a jwt with the code and the code is saved in the database with the authorized user id and redirected
## get accessToken by the code from the backchannel
 http://localhost:8080/login/oauth/access_token?client_id=84bd6d67-892d-4ebb-ab8f-361504388042&client_secret=TfF1Zz96ExoqYu3Ka5lP&code=code
 the above url with get the code from the jwt and check the database for the corresponding userid and make a jwt with the user and send as JSON
## get userinfo by the accessToken
 POST
 Header :- Authentication : Bearer access_token
 URL := http://localhost:8080/login/oauth/user
 the above url will get the userid from the jwt get the user from the database and send the user as JSON


 ## TODO

- scopes and openid
- mysql implimentations
- Redirect to unknown error in line 85 of oAuth.go
- Validations in the backend
- Postgres connection error
- Add a login time table to check for bruteforce
- CSRF Token
- helmet
- better ratelimits
- css cleanups