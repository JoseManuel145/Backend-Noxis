package domain

type Kit struct {
	Clave    string   `json:"clave"`
	Sensores []string `json:"sensores"`
	Username string   `json:"username"`
	Userfk   int      `json:"userfk,omitempty"`
	Status   bool     `json:"status"` //true = "canjeado" || false = "no canjeado"
}
