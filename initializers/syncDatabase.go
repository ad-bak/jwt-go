package initializers

import "jwt/go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
