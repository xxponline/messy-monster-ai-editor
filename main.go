package main

import (
	"github.com/gin-gonic/gin"
	"messy-monster-ai-editor/asset_content"
	"messy-monster-ai-editor/asset_organization"
	"messy-monster-ai-editor/db"
)

func main() {
	db.ServerDatabase = &db.SqliteDataBase{}
	errCode, errMsg := db.ServerDatabase.Initialize("./db/db.sqlite")
	if errCode != 0 {
		panic(errMsg)
	}

	//var items []common.SolutionInfoItem
	//errCode, errMsg, items = solutionMgr.ListSolutions()
	//fmt.Printf("===>>> %d \n", errCode)
	//
	//fmt.Printf("%v", items)

	//uuid := uuid.New()
	//fmt.Println(uuid.String())

	r := gin.Default()
	APIRout := r.Group("API")
	asset_organization.InitializeAssetManagement(APIRout.Group("AssetManagement"))
	asset_content.InitializeAssetManagement(APIRout.Group("AssetContentModifier"))
	r.Run("localhost:8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//type FOO struct {
//	URL      string `json:"url" binding:"required"`
//	Solution string `json:"solutionName" binding:"required"`
//}
//
//func main() {
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		c.String(200, "pong")
//	})
//	r.POST("/foo", func(c *gin.Context) {
//		var url FOO
//		c.BindJSON(&url)
//		fmt.Printf("URL to store: %v\n", url)
//	})
//	r.Run(":8080") // listen an
//}
