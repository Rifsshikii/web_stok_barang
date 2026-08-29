package models

type Admin struct {
	IDAdmin  int    `json:"id_admin"`
	Nama     string `json:"nama"`
	NomorHP  string `json:"nomor_hp"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
