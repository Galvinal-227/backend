package main

import (
	"log"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/database"
	"pkl-attendance/backend/models"
	"pkl-attendance/backend/utils"
)

func main() {
	cfg, _ := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	_ = database.Migrate(db)

	var existing models.User
	if err := db.Where("username = ?", "siswa1").First(&existing).Error; err == nil {
		log.Println("siswa1 sudah ada")
		return
	}
	hash, _ := utils.HashPassword("Siswa123!")
	u := &models.User{Username: "siswa1", Password: hash, Name: "Siswa Uji", Role: models.RoleSiswa, IsActive: true}
	if err := db.Create(u).Error; err != nil {
		log.Fatal(err)
	}

	start := models.DateOnlyFromTime(utils.JakartaNow())
	end := models.DateOnlyFromTime(utils.JakartaNow().AddDate(0, 3, 0))
	s := &models.Student{UserID: u.ID, NIS: "12345", Class: "XII RPL 1", Major: "RPL", PKLStartDate: &start, PKLEndDate: &end}
	if err := db.Create(s).Error; err != nil {
		log.Fatal(err)
	}
	log.Println("siswa1 dibuat (password: Siswa123!)")
}