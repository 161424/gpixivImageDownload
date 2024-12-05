package model

import "time"

var Cacher = Cache{}

type Cache struct {
	Time time.Time `json:"time"`
	*Common
	*Ranks
	*Author
}

func init() {
	Cacher.Time = time.Now()
}

func (c *Cache) Update() {
	c.Time = time.Now()

}
