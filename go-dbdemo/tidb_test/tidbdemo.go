package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open("mysql", "gotest:panenming@tcp(10.39.35.38:4000)/gotest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		panic("数据库连接失败！")
	}

	defer db.Close()

	//	db.DropTable(&Product{})
	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "L1212", Price: 1000})
	// 读取
	var product Product
	db.First(&product, 1) // 查询id为1的product
	fmt.Println("code=", product.Code)
	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	fmt.Println("code=", product.Code)
	// 更新 - 更新product的price为2000
	db.Model(&product).Update("Price", 2000)
	fmt.Println("price=", product.Price)
	// 删除 - 删除product
	db.Delete(&product)

}
