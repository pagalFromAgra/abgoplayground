package uploadlogs

type Payload struct {
	DateTime    string `json:"time" binding:"requred"`
	Company     string `json:"company" binding:"required"`
	Device      string `json:"device" binding:"required"`
	Event       string `json:"event" binding:"required"`
	Vbat        int64  `json:"vbat" binding:"required"`
	Temperature int64  `json:"temperature" binding:"required"`
	Value       int64  `json:"value" binding:"required"` // receive as a string, otherwise value = 0 will be omitted from json
	Info        string `json:"info" binding:"required"`
}
