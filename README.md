# Rental Property API

## Setup and Run

```bash
bee run
```
The API will run at [http://localhost:8080](http://localhost:8080)

## Swagger UI

Open the Swagger UI at: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

## Some Sample curl commands

```bash
curl "http://localhost:8080/v1/properties" > res.json
```

```bash
curl "http://localhost:8080/v1/properties?min_price=100&max_price=500&min_star_rating=4&limit=10" > res.json
```

```bash
curl "http://localhost:8080/v1/properties?amenities=Breakfast%20Included&limit=10" > res.json
```

```bash
curl "http://localhost:8080/v1/properties/1" > res.json
```

## Testing

Testing can be performed using:
```bash
go test ./services -v
```

The result will be:

```
?       rental-property-api     [no test files]
?       rental-property-api/controllers [no test files]
?       rental-property-api/models      [no test files]
?       rental-property-api/routers     [no test files]
=== RUN   TestGetResponses_Transform
--- PASS: TestGetResponses_Transform (0.00s)
=== RUN   TestGetResponses_FilterAND
=== RUN   TestGetResponses_FilterAND/minimum_price_boundary
=== RUN   TestGetResponses_FilterAND/maximum_price_boundary
=== RUN   TestGetResponses_FilterAND/price_range
=== RUN   TestGetResponses_FilterAND/minimum_star_rating
=== RUN   TestGetResponses_FilterAND/minimum_review_score
=== RUN   TestGetResponses_FilterAND/minimum_reviews
=== RUN   TestGetResponses_FilterAND/published_true
=== RUN   TestGetResponses_FilterAND/published_false
=== RUN   TestGetResponses_FilterAND/property_type
=== RUN   TestGetResponses_FilterAND/feed
=== RUN   TestGetResponses_FilterAND/minimum_bedroom
=== RUN   TestGetResponses_FilterAND/multiple_scalar_filters_use_AND
=== RUN   TestGetResponses_FilterAND/no_property_satisfies_filter
--- PASS: TestGetResponses_FilterAND (0.00s)
    --- PASS: TestGetResponses_FilterAND/minimum_price_boundary (0.00s)
    --- PASS: TestGetResponses_FilterAND/maximum_price_boundary (0.00s)
    --- PASS: TestGetResponses_FilterAND/price_range (0.00s)
    --- PASS: TestGetResponses_FilterAND/minimum_star_rating (0.00s)
    --- PASS: TestGetResponses_FilterAND/minimum_review_score (0.00s)
    --- PASS: TestGetResponses_FilterAND/minimum_reviews (0.00s)
    --- PASS: TestGetResponses_FilterAND/published_true (0.00s)
    --- PASS: TestGetResponses_FilterAND/published_false (0.00s)
    --- PASS: TestGetResponses_FilterAND/property_type (0.00s)
    --- PASS: TestGetResponses_FilterAND/feed (0.00s)
    --- PASS: TestGetResponses_FilterAND/minimum_bedroom (0.00s)
    --- PASS: TestGetResponses_FilterAND/multiple_scalar_filters_use_AND (0.00s)
    --- PASS: TestGetResponses_FilterAND/no_property_satisfies_filter (0.00s)
=== RUN   TestGetResponses_FilterAmenitiesOR
=== RUN   TestGetResponses_FilterAmenitiesOR/first_requested_amenity_matches
=== RUN   TestGetResponses_FilterAmenitiesOR/second_requested_amenity_matches
=== RUN   TestGetResponses_FilterAmenitiesOR/none_match
=== RUN   TestGetResponses_FilterAmenitiesOR/case_sensitive
--- PASS: TestGetResponses_FilterAmenitiesOR (0.00s)
    --- PASS: TestGetResponses_FilterAmenitiesOR/first_requested_amenity_matches (0.00s)
    --- PASS: TestGetResponses_FilterAmenitiesOR/second_requested_amenity_matches (0.00s)
    --- PASS: TestGetResponses_FilterAmenitiesOR/none_match (0.00s)
    --- PASS: TestGetResponses_FilterAmenitiesOR/case_sensitive (0.00s)
=== RUN   TestGetResponses_FilterCombined
--- PASS: TestGetResponses_FilterCombined (0.00s)
=== RUN   TestGetResponses_Limit
=== RUN   TestGetResponses_Limit/limit_one
=== RUN   TestGetResponses_Limit/limit_two
=== RUN   TestGetResponses_Limit/limit_greater_than_result
--- PASS: TestGetResponses_Limit (0.00s)
    --- PASS: TestGetResponses_Limit/limit_one (0.00s)
    --- PASS: TestGetResponses_Limit/limit_two (0.00s)
    --- PASS: TestGetResponses_Limit/limit_greater_than_result (0.00s)
=== RUN   TestGetResponses_EmptyResult
--- PASS: TestGetResponses_EmptyResult (0.00s)
=== RUN   TestGetResponseByID
=== RUN   TestGetResponseByID/found
=== RUN   TestGetResponseByID/not_found
--- PASS: TestGetResponseByID (0.00s)
    --- PASS: TestGetResponseByID/found (0.00s)
    --- PASS: TestGetResponseByID/not_found (0.00s)
=== RUN   TestLoadJSON
=== RUN   TestLoadJSON/valid_JSON
=== RUN   TestLoadJSON/missing_file
=== RUN   TestLoadJSON/invalid_JSON
--- PASS: TestLoadJSON (0.00s)
    --- PASS: TestLoadJSON/valid_JSON (0.00s)
    --- PASS: TestLoadJSON/missing_file (0.00s)
    --- PASS: TestLoadJSON/invalid_JSON (0.00s)
PASS
ok      rental-property-api/services    (cached)
testing: warning: no tests to run
PASS
ok      rental-property-api/tests       (cached) [no tests to run]
```

and, to see the test coverage:

```bash
go test ./... -cover
```

which will show the following output:

```
        rental-property-api             coverage: 0.0% of statements
        rental-property-api/controllers         coverage: 0.0% of statements
?       rental-property-api/models      [no test files]
        rental-property-api/routers             coverage: 0.0% of statements
ok      rental-property-api/services    (cached)        coverage: 98.5% of statements
ok      rental-property-api/tests       (cached)        coverage: [no statements] [no tests to run]
```

that says the Test Coverage for services/ is 98.5%