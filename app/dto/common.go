package dto

type HandleError struct {
	Status int   `json:"status"`
	Errors error `json:"errors"`
}
