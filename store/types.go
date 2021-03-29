package store

type Pilot struct {
	Cid         string  `json:"cid"`
	Server      string  `json:"server"`
	Rating      string  `json:"rating"`
	RealName    string  `json:"realName"`
	Callsign    string  `json:"callsign"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Heading     float64 `json:"heading"`
	Speed       float64 `json:"speed"`
	Transponder string  `json:"transponder"`
	Type        string  `json:"-"`
}
