# Rental Property API

## Setup and Run

```bash
bee run
```
The API will run at [http://localhost:8080](http://localhost:8080)

## Sample curl commands

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
