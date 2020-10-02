# gen-test

## How to run:

> `cd gen-test`
> `docker-compose up`


## Try accessing following APIs using Postman (or other REST clients):

1. http://localhost:5000/post-data

`Request data:
{
    "Title": "Hello"
}`

Sample Response:
`{
    "Title": "Hello",
    "UUID4": "eb8d291c-3f86-462a-9ad3-6958f26225c8",
    "Timestamp": "2020-10-02T12:11:36.219666Z"
}`

2. http://localhost:5000/get-data/eb8d291c-3f86-462a-9ad3-6958f26225c8
`{
    "Title": "Hello",
    "UUID4": "eb8d291c-3f86-462a-9ad3-6958f26225c8",
    "Timestamp": "2020-10-02T12:11:36.219666Z"
}`


## Notes:

1. Used Docker version 17.12.0-ce
2. I have pushed .env file to git so this can be tested directly, ideally this file will be added in .gitignore
