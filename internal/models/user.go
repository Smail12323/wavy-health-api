package models

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Balance      float64
	StocksOwned  int
}
