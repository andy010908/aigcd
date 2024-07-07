package mj

import (
	log "aigcd/core/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	core "aigcd/core/mysql"

	"cloud.google.com/go/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/golangFame/imageslicer"
	"go.uber.org/zap"
)

func DiscordBot() {
	log.Info("DiscordBot starting ...")
	// Replace with your bot token
	botToken := "MTA4ODM4ODI1MTQ4MTU1NDk1NA.GvKKWi.6j6ZRxf1UY2QWo2hjQ1odIDWwRum8B8UeEGgnk"

	// Replace with the channel ID you want to listen to
	channelID := "1088383794702196829"

	// Create a new Discord session using the bot token
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Add a message create handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log.Info("DiscordBot", zap.String("ID", m.ID),
			zap.String("ChannelID", m.ChannelID),
			zap.String("Content", m.Content))

		if m.ChannelID != channelID {
			return
		}

		// Check if the message has any attachments
		if len(m.Attachments) > 0 {
			// Loop through each attachment
			for _, attachment := range m.Attachments {
				// Check if the attachment is an image
				log.Info("DiscordBot", zap.String("attachment.URL", attachment.URL),
					zap.String("attachment.Filename", attachment.Filename),
					zap.String("attachment.ContentType", attachment.ContentType))

				//if strings.HasPrefix(attachment.ContentType, "image/") {
				// Download the image
				err := downloadImage(attachment.URL, attachment.Filename)
				if err != nil {
					log.Error("Error downloading image: " + err.Error())
				} else {
					log.Info("Image downloaded successfully: " + attachment.Filename)
					go splitImage(attachment.Filename)
				}
				//}
			}
		}
	})

	// Open the Discord session
	err = dg.Open()
	if err != nil {
		log.Error("Error opening Discord session: " + err.Error())
		return
	}

	// Wait for a signal to stop the bot
	log.Info("DiscordBot is running ...")
	<-make(chan struct{})
	dg.Close()
	log.Info("DiscordBot stopped")
}

// Function to download an image from a URL and save it to disk
func downloadImage(url string, filename string) error {
	// Create a new HTTP request to download the image
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a new file to save the image
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the image from the HTTP response to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

const cloudStorageURL string = "https://storage.googleapis.com/aigcd/"

func splitImage(path string) {
	img, _ := imageslicer.GetImageFromPath(path)
	if img == nil {
		log.Error("splitImage:invalid image url or image format not supported!")
	}
	grid := imageslicer.Grid{2, 2} //rows,columns
	tiles := imageslicer.Slice(img, grid)
	expectedTiles := int(grid[0] * grid[1])
	if len(tiles) != expectedTiles {
		log.Error("splitImage: slice error")
	}
	str := strings.Split(path, ".")
	var fileName string
	for k, v := range tiles {
		fileName = fmt.Sprintf("%s_%d.%s", str[0], k, str[1])
		err := imageslicer.Save(v, fileName)
		if err != nil {
			log.Error("splitImage:" + err.Error())
			break
		}
		if k == 0 { //only upload one image
			res := UploadFile("aigcd", fileName)
			if res == true {
				addRecord(cloudStorageURL + fileName)
				//delete file
				err := os.Remove(fileName)
				if err != nil {
					log.Error("splitImage: " + err.Error())
				}
				err = os.Remove(path)
				if err != nil {
					log.Error("splitImage: " + err.Error())
				}
			}
			break
		}
	}

}

func UploadFile(bucket, object string) bool {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error("storage.NewClient: " + err.Error())
		return false
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(object)
	if err != nil {
		log.Error("os.Open: " + err.Error())
		return false
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		log.Error("io.Copy: " + err.Error())
		return false
	}
	if err := wc.Close(); err != nil {
		log.Error("Writer.Close: " + err.Error())
		return false
	}
	return true
}

func addRecord(url string) error {
	image := core.Collections{
		CreatedAt: time.Now().Unix(),
		ImageURL:  url,
	}
	result := core.DB.Create(&image)
	return result.Error
	//fmt.Println(result.Error)        // nil
	//fmt.Println(result.RowsAffected)
}
