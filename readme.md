
```markdown
# Football League Backend API

## 1. Project Overview

This project is a backend API developed with **GoLang** for managing a football league system.

The system allows users to:

- Create and list teams
- Generate league matches
- Update match results
- Calculate league standings
- Predict championship probabilities using simulation

The project uses **PostgreSQL** as the database and runs with **Docker / Docker Compose**.

---

## 2. Technologies Used

- GoLang
- PostgreSQL
- Docker
- Docker Compose
- REST API
- Postman

---

## 3. Project Structure

```text
.
├── cmd/
│   └── main.go
├── internal/
│   ├── database/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── service/
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md

```

#### Test the deployed project from this base url

http://13.62.101.164:8080/

*(Url may not be active at later times.)*

---

## 4. API Endpoints

### Frontend page for testing api

`/`

---

## Team Endpoints

### Create Team

**POST** `/team/create`

Creates a new team and returns the created team data.

The `strength` value must be a float between `0` and `1`.

**Headers**

```http
Content-Type: application/json

```

**Request Body**

```json
{
  "name": "Çorum spor",
  "strength": 0.5
}

```

**Example Response**

```json
{
  "id": 1,
  "name": "Çorum spor",
  "strength": 0.5
}

```

---

### Get All Teams

**GET** `/team/get_all`

Returns all teams as JSON.

**Example Response**

```json
[
  {
    "id": 4,
    "name": "Galatasaray",
    "strength": 0.8
  },
  {
    "id": 5,
    "name": "Fener",
    "strength": 0.7
  }
]

```

---

### Get Team By Id

**GET** `/team/get_by_id/{id}`

Returns a specific team by ID as JSON.

**Example Response**

```json
{
  "id": 4,
  "name": "Galatasaray",
  "strength": 0.8
}

```

---

### Delete All Teams

**DELETE** `/team/delete_all`

Deletes all teams and related match records.

---

## Match Endpoints

### Generate Matches

**POST** `/match/generate`

Generates the match schedule randomly according to the given total week number.

**Headers**

```http
Content-Type: application/json

```

**Request Body**

```json
{
  "total_week": 3
}

```

**Example Response**

```json
[
  {
    "id": 13,
    "week_number": 1,
    "home_team_id": 4,
    "away_team_id": 5,
    "home_score": 0,
    "away_score": 0,
    "is_played": false
  }
]

```

---

### Get All Matches

**GET** `/match/get_all`

Returns all match records.

---

### Update Match Results

**PUT** `/match/update/{match_id}`

Update Match Results and returns the resulting match status.

**Headers**

```http
Content-Type: application/json

```

**Request Body**

```json
{
  "home_score": 2,
  "away_score": 1,
  "is_played": true
}

```

**Example Response**

```json
{
  "id": 13,
  "week_number": 1,
  "home_team_id": 4,
  "away_team_id": 5,
  "home_score": 2,
  "away_score": 1,
  "is_played": true
}

```

---

### Delete All Matches

**DELETE** `/match/delete_all`

Deletes all match records, including both played and unplayed matches.

---

## League Endpoints

### Play Current Week

**POST** `/league/play_current_week`

Simulates the current week matches and returns the played match results.

**Request Body**

```json
{
  "unpred_coef": 0.01
}

```

---

### Play All Weeks

**POST** `/league/play_all_weeks`

Simulates all unplayed matches until the end of the season.

---

### Get League Table

**GET** `/league/get_table`

Returns the current league table in ordered form.

---

### Simulate Championship Probabilities

**POST** `/league/simulate`

Runs Monte Carlo simulations to estimate each team's championship probability.

**Request Body**

```json
{
  "unpred_coef": 0.01
}

```
