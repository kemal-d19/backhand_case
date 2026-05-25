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

(Url may not be active later times.)

```markdown
## 4. API Endpoints

This section describes the available API endpoints for teams, matches, and league simulation.

---
### Frontend page for testing api 

`/`

---

## Team Endpoints

### Create Team

**POST** `/team/create`

Creates a new team and returns the created team data.

The `strength` value must be a float between `0` and `1`.

**Headers**

http
Content-Type: application/json
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

**GET** `/team/get_by_id/<id>`

Returns all teams as JSON.

**Example Response**

```json
[
  {
    "id": 4,
    "name": "Galatasaray",
    "strength": 0.8
  },
]
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

All generated matches are created with `is_played = false`.

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
  },
  {
    "id": 17,
    "week_number": 1,
    "home_team_id": 5,
    "away_team_id": 6,
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
  },
  {
    "id": 17,
    "week_number": 1,
    "home_team_id": 5,
    "away_team_id": 6,
    "home_score": 0,
    "away_score": 0,
    "is_played": false
  }
]
```
---

### Update Match Results
**GET** `/match/get_all<match_id>`

Update Match Results and returns resulting match status..

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
[
  {
    "id": 13,
    "week_number": 1,
    "home_team_id": 4,
    "away_team_id": 5,
    "home_score": 2,
    "away_score": 1,
    "is_played": true
  }
]
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

The `unpred_coef` value must be between `0` and `1`.

This value controls the unpredictability of the match results. In simple terms, a higher value increases randomness.

**Headers**

```http
Content-Type: application/json
```

**Request Body**

```json
{
  "unpred_coef": 0.01
}
```

**Example Response**

```json
[
  {
    "id": 19,
    "week_number": 1,
    "home_team_id": 4,
    "away_team_id": 5,
    "home_score": 4,
    "away_score": 0,
    "is_played": true
  },
  {
    "id": 21,
    "week_number": 1,
    "home_team_id": 4,
    "away_team_id": 6,
    "home_score": 1,
    "away_score": 0,
    "is_played": true
  }
]
```

---

### Play All Weeks

**POST** `/league/play_all_weeks`

Simulates all unplayed matches until the end of the season.

The `unpred_coef` value controls the randomness of the simulated results.

**Headers**

```http
Content-Type: application/json
```

**Request Body**

```json
{
  "unpred_coef": 0.01
}
```

**Response**

Returns all updated match records after simulation.

---

### Get League Table

**GET** `/league/get_table`

Returns the current league table in ordered form.

Teams are ordered according to their league performance.

**Example Response**

```json
[
  {
    "team": "Galatasaray",
    "pts": 6,
    "played": 2,
    "won": 2,
    "drawn": 0,
    "lost": 0,
    "gd": 1
  },
  {
    "team": "Çorum spor",
    "pts": 0,
    "played": 1,
    "won": 0,
    "drawn": 0,
    "lost": 1,
    "gd": -1
  }
]
```

---

### Simulate Championship Probabilities

**POST** `/league/simulate`

Runs Monte Carlo simulations to estimate each team's championship probability based on the currently played matches.

The returned `prediction` value represents the estimated championship probability as a percentage.

**Headers**

```http
Content-Type: application/json
```

**Request Body**

```json
{
  "unpred_coef": 0.01
}
```

**Example Response**

```json
[
  {
    "team": "Galatasaray",
    "prediction": 81
  },
  {
    "team": "Fener",
    "prediction": 13.2
  },
  {
    "team": "Çorum spor",
    "prediction": 5.8
  }
]
```
```
