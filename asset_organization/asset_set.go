package asset_organization

import (
	"github.com/gin-gonic/gin"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
	"net/http"
)

type ListAssetSetReq struct {
	SolutionId string `json:"solutionId" binding:"required"`
}

type CreateAssetSetReq struct {
	SolutionId   string `json:"solutionId" binding:"required"`
	AssetSetName string `json:"assetSetName" binding:"required"`
}

func ListAssetSetsAPI(context *gin.Context) {
	var req ListAssetSetReq
	context.BindJSON(&req)

	var assetSetInfos []common.AssetSetInfoItem

	errCode, errMsg, assetSetMgr := db.ServerDatabase.GetAssetSetManager(false)
	if errCode == common.Success {
		defer assetSetMgr.Release()
		errCode, errMsg, assetSetInfos = assetSetMgr.ListAssetSets(req.SolutionId)
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"assetSets":  assetSetInfos,
	})
}

func CreateAssetSetAPI(context *gin.Context) {
	var req CreateAssetSetReq
	context.BindJSON(&req)

	var assetSetInfos []common.AssetSetInfoItem

	errCode, errMsg, assetSetMgr := db.ServerDatabase.GetAssetSetManager(true)
	if errCode == common.Success {
		defer assetSetMgr.Release()
		errCode, errMsg = assetSetMgr.CreateAssetSet(req.SolutionId, req.AssetSetName)
		if errCode == common.Success {
			errCode, errMsg, assetSetInfos = assetSetMgr.ListAssetSets(req.SolutionId)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"assetSets":  assetSetInfos,
	})

}
