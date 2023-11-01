package main

import (
	"cache/internal/model/information"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/goim?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// db.AutoMigrate(&information.Information{})

	Info := &information.Information{}
	var t = "Tom"
	result := db.Where("k = ?", t).First(Info)
	fmt.Print(result, Info)
	// db.AutoMigrate(&file.File{})
}
