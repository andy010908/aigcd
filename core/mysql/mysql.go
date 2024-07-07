package core

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Collections struct {
	CreatedAt int64  `json:"created_at" `
	ImageURL  string `json:"image_url" `
}

const (
	UserName     string = "aigc"
	Password     string = "laputa0314"
	Addr         string = "35.197.95.208"
	Port         int    = 3306
	Database     string = "aigc"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

var DB *gorm.DB

func init() {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", UserName, Password, Addr, Port, Database)
	var err error
	DB, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}

	sql, err := DB.DB()
	if err != nil {
		fmt.Println("sql failed:", err)
		return
	}
	sql.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	sql.SetMaxIdleConns(MaxIdleConns)
	sql.SetMaxOpenConns(MaxOpenConns)

	/*
		image := Collections{
			CreatedAt: time.Now().Unix(),
			ImageURL:  "https://storage.cloud.google.com/aigcd/harrison0617_a_picture_from_a_brewery_coffee_shop_kitchen_outdo_1ad5d2b0-938b-465b-9154-85bee66fd181_0.png",
		}
		result := DB.Create(&image)
		fmt.Println(result.Error) // nil
		fmt.Println(result.RowsAffected)
	*/

}
