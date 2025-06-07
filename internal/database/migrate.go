package database

import (
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	"log"
)

func UpdateMigration() {
	// Step 1: Check if the Message table exists
	if !DB.Migrator().HasTable(&echoModel.Message{}) {
		log.Println("Message table does not exist, skipping data migration.")
		return
	}

	// Step 2: If the Message table exists, proceed with migration
	var messages []echoModel.Message
	if err := DB.Find(&messages).Error; err != nil {
		log.Fatal("failed to fetch messages", err)
	}

	// Step 3: Insert data into the Echo table
	for _, msg := range messages {
		echo := echoModel.Echo{
			Content:       msg.Content,
			Username:      msg.Username,
			Private:       msg.Private,
			UserID:        msg.UserID,
			Extension:     msg.Extension,
			ExtensionType: msg.ExtensionType,
			CreatedAt:     msg.CreatedAt,
		}

		// Create Echo record
		if err := DB.Create(&echo).Error; err != nil {
			log.Printf("failed to create echo for message ID %d: %v", msg.ID, err)
		}

		// Step 4: Migrate associated images
		for _, img := range msg.Images {
			image := echoModel.Image{
				MessageID:   echo.ID, // Associate the image with the new echo
				ImageURL:    img.ImageURL,
				ImageSource: img.ImageSource,
			}
			if err := DB.Create(&image).Error; err != nil {
				log.Printf("failed to create image for message ID %d: %v", msg.ID, err)
			}
		}
	}

	log.Println("Data migration complete")
}
