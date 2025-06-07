package database

import (
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	"log"
)

func UpdateMigration() {
	// Step 1: Check if the Message table exists
	if DB.Migrator().HasTable(&echoModel.Message{}) {
		log.Println("Message table exists, indicating user is using an old version. Proceeding with data migration.")

		// Step 2: Fetch all messages from the old Message table
		var messages []echoModel.Message
		if err := DB.Find(&messages).Error; err != nil {
			log.Fatal("failed to fetch messages", err)
		}

		// Step 3: Insert data into the Echo table, avoid duplicates
		for _, msg := range messages {
			// Check if the message has already been migrated to Echo
			var existingEcho echoModel.Echo
			if err := DB.Where("id = ?", msg.ID).First(&existingEcho).Error; err == nil {
				// If record exists, skip this message
				//log.Printf("Echo record for message ID %d already exists, skipping.", msg.ID)
				continue
			}

			// If not, insert the Echo record
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
				continue
			}

			// Step 4: Migrate associated images, avoid duplicates
			for _, img := range msg.Images {
				var existingImage echoModel.Image
				if err := DB.Where("message_id = ? AND image_url = ?", echo.ID, img.ImageURL).First(&existingImage).Error; err == nil {
					// If image already exists, skip it
					log.Printf("Image for message ID %d already exists, skipping.", msg.ID)
					continue
				}

				// Create Image record
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
	} else {
		// If the Message table doesn't exist, skip migration
		log.Println("Message table does not exist, skipping data migration. User is on the latest version.")
	}
}
