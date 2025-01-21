# Receipt Processor Challenge

## Requirements
As outlined in the go.mod, I used go v1.20 and the gorilla/mux library

`go mod tidy`  
`go run receipt-processor-challenge.go`  
  
It will run on port 8080

## Main
Initialize a map of string to int to keep track of the points for each receipt
Set up a router for the endpoints

## Endpoints
### POST /receipts/process  

- Decode the json body into a "Receipt" struct matching the schema in the api.yml  
  - Return "The receipt is invalid." if receipt can't be decoded / doesn't match schema  
- Create an id by concatenating retailer, date, and time  
  - Remove non alphanuemeric characters from retailer name so it plays nicely with the url  
- Calculate the points according to the rules  
  - I'm not a large language model  
- Store the points in the map  
- Return the json {"id": "{id_val}"}  

### GET /receipts/{id}/points  

- Get the id string from the url parameter  
- Check the map  
  - Return "No receipt found for that ID." if id not present  
- Return the json {"points": "{point_val}"} with the value calculated previously and stored in the map

## Testing
  
- I tested each endpoint using Postman
- An id was returned for the example receipts
- The correct point value was stored in the map and returned upon Get request to /receipts/{id}/points
- Get requests to /receipts/{id}/points for id values not present return a 404 with error message: "No receipt found for that ID."
