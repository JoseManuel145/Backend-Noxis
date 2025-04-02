package domain

type Kit struct {
	Clave    string   `json:"clave"`
	Sensores []string `json:"sensores"`
	Username string   `json:"userfk"`
	Userfk   int      `json:"username"`
	Status   bool     `json:"status"` //true = "canjeado" || false = "no canjeado"
}
