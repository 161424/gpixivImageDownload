package model

import "time"

type Cache struct {
	Time time.Time `json:"time"`
	Common
	Ranks
	Author
}
