package model

type Player struct {
	Name       string `json:"name"`
	Number     string `json:"number"`
	Position   string `json:"position"`
	Age        string `json:"age"`
	Birthplace string `json:"birthplace"`
	ImagePath  string `json:"image"`
	PlayerURL  string `json:"url"`
}

