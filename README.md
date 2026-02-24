# shoppinglist

## Demo

- [Screenshot](./demo/demo.png)
- [Video](./demo/demo.webm)

## Creating the repo

- `npm create vite@latest shoppinglist -- --template react-ts`
- `npm install antd --save`
- `npm i @tanstack/react-query`
- `go mod init github.com/tomp-work/shoppinglist`
- `go get github.com/labstack/echo/v5`
- `go get github.com/stretchr/testify`
- `go get github.com/resend/resend-go/v2`
- `go get github.com/joho/godotenv`
- `go mod tidy`
- `npm install @ant-design/icons@6.x --save`

## Running

Backend:

- `go run cmd/server/main.go`

Frontend:

- `npm run dev`

## Test cURL

cURL command to test sending a shopping list by email:

```sh
curl -H "Content-Type: application/json" \
  -X POST \
  -d '{"emailAddress":"valid.email@example.com"}' \
  http://localhost:1323/list/send
```

## TODO

- [ ] Consider using UUIDs
- [ ] Improve FE status/error handling

## DONE

- [x] Dev env and tech stack setup
- [x] Story1: get shopping list items: BE
- [x] Story1: get shopping list items: FE
- [x] Story2: add shopping list item: BE
- [x] Story2: add shopping list item: FE
- [x] Story3: delete shopping list item: BE
- [x] Story3: delete shopping list item: FE
- [x] Story4: Cross off shopping list item: BE
- [x] Story4: Cross off shopping list item: FE
- [x] Story5: Persist the list
- [x] Story6: change order of shopping list items: BE
- [x] Story6: change order of shopping list items: FE
- [x] Story7: shopping list total price: BE
- [x] Story7: shopping list total price: FE
- [x] Fix generateId bug (bug: IDs were being reused after delete).
- [x] Story8: spending limit: BE
- [x] Refactor: total calculation is done on create/delete and persisted in ListDetails.
- [x] Story8: spending limit: FE
- [x] Remove quantity from ListItem
- [x] Story9: send shopping list to email address: BE
- [x] Story9: send shopping list to email address: FE

