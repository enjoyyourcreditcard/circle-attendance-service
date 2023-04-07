package domain

type WebResponse struct {
	Status 		int				`json:"status"`
	Data 		interface{}		`json:"data"`
	Message		string			`json:"message"`
}