package views

type RequestAlamat struct {
	Idalamat     int    `gorm:"primaryKey;autoincrement" json:"idalamat"`
	IdUser       int    `json:"iduser" binding:"required"`
	AlamatDetail string `json:"alamatdetail" binding:"required"`
	Provinsi     string `json:"provinsi" binding:"required"`
	Kabupaten    string `json:"kabupaten" binding:"required"`
	Kodepos      string `json:"kodepos" binding:"required"`
}
