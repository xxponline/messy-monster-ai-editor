package asset_organization

import (
	"github.com/gin-gonic/gin"
	"messy-monster-ai-editor/asset_content/content_modifier"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
	"net/http"
)

type CreateAssetReq struct {
	AssetSetId string `json:"assetSetId" binding:"required"`
	AssetType  string `json:"assetType" binding:"required"`
	AssetName  string `json:"assetName" binding:"required"`
}

type ListAssetsReq struct {
	AssetSetId string `json:"assetSetId" binding:"required"`
}

type ListAssetsByMultipleSetsReq struct {
	AssetSetIds []string `json:"assetSetIds" binding:"required"`
}

type ReadAssetReq struct {
	AssetId string `json:"assetId" binding:"required"`
}

func CreateAssetAPI(context *gin.Context) {
	var req CreateAssetReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	var errCode common.ErrorCode
	var errMsg string
	var initialContent string

	//assetType Check
	switch req.AssetType {
	case "BehaviourTree":
		errCode, errMsg, initialContent = content_modifier.BehaviourTreeCreateEmptyContent()
		if errCode != common.Success {
			context.JSON(http.StatusOK, gin.H{
				"errCode":    errCode,
				"errMessage": errMsg,
			})
			return
		}
	case "BlackBoard":
	default:
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.InvalidAssetType,
			"errMessage": common.InvalidAssetType.GetMsgFormat(req.AssetType),
		})
		return
	}

	var assetItems []common.AssetSummaryInfoItem
	var assetMgr db.IAssetManager

	errCode, errMsg, assetMgr = db.ServerDatabase.GetAssetManager(true)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}
	defer assetMgr.Release()

	errCode, errMsg, createAssetId := assetMgr.CreateAsset(req.AssetSetId, req.AssetType, req.AssetName, initialContent)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	errCode, errMsg, assetItems = assetMgr.ListAssets([]string{req.AssetSetId})

	context.JSON(http.StatusOK, gin.H{
		"errCode":           errCode,
		"errMessage":        errMsg,
		"assetSummaryInfos": assetItems,
		"newAssetId":        createAssetId,
	})
}

func ListAssetsAPI(context *gin.Context) {
	var req ListAssetsReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	var assetItems []common.AssetSummaryInfoItem

	errCode, errMsg, assetMgr := db.ServerDatabase.GetAssetManager(false)
	if errCode == common.Success {
		defer assetMgr.Release()
		errCode, errMsg, assetItems = assetMgr.ListAssets([]string{req.AssetSetId})
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":           errCode,
		"errMessage":        errMsg,
		"assetSummaryInfos": assetItems,
	})
}

func ListAssetsByMultipleAssetSetsAPI(context *gin.Context) {
	var req ListAssetsByMultipleSetsReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	var assetItems []common.AssetSummaryInfoItem

	errCode, errMsg, assetMgr := db.ServerDatabase.GetAssetManager(false)
	if errCode == common.Success {
		defer assetMgr.Release()
		errCode, errMsg, assetItems = assetMgr.ListAssets(req.AssetSetIds)
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":           errCode,
		"errMessage":        errMsg,
		"assetSummaryInfos": assetItems,
	})
}

func ReadAssetAPI(context *gin.Context) {
	var req ReadAssetReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	var assetDetail *common.AssetDetailInfo

	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(req.AssetId, false)
	if errCode == common.Success {
		defer assetDoc.Release()
		errCode, errMsg, assetDetail = assetDoc.ReadAsset()
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":       errCode,
		"errMessage":    errMsg,
		"assetDocument": assetDetail,
	})
}
