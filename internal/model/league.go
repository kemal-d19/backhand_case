package model

type LeagueTableRow struct {
	Team   string `json:"team" db:"Team"`
	Points int    `json:"pts" db:"points"`
	Played int    `json:"played" db:"played"`
	Won    int    `json:"won" db:"won"`
	Drawn  int    `json:"drawn" db:"drawn"`
	Lost   int    `json:"lost" db:"lost"`
	GD     int    `json:"gd" db:"gd"`
}

type PredictionTableRow struct {
	Team       string  `json:"team" db:"team"`
	Prediction float32 `json:"prediction" db:"prediction"`
}
