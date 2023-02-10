#!/bin/zsh

# register redirect
curl -v -X POST localhost:8080/google -d 'https://google.com'

# access redirect
curl -v http://localhost:8080/google

# unregister redirect
curl -v -X DELETE localhost:8080/google

# list redirects
curl -v http://localhost:8080/list 