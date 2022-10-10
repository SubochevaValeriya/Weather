package weather

import "time"

type Subscription struct {
	Id               int       `json:"id" db:"id"`
	City             string    `json:"city" db:"city"`
	SubscriptionDate time.Time `json:"subscription_date" db:"subscription_date"`
}

type Weather struct {
	City        string    `json:"name" db:"city"`
	Temperature int       `json:"temperature" db:"temperature"`
	Date        time.Time `json:"date" db:"date"`
}

type City struct {
	City string `json:"name" db:"city"`
}
