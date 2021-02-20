package shared

type ForwardResponse struct {
	ChiselId string `json:"chiselId"`
	ChiselUsername string `json:"username"`
	ChiselPassword string `json:"password"`
	AllowedPort int `json:"allowedPort"`
	TunnelUID string `json:"tunnelUid"`
}
