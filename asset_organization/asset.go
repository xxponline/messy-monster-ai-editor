package asset_organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xxponline/messy-monster-ai-editor/asset_content/content_modifier"
	"github.com/xxponline/messy-monster-ai-editor/common"
	"github.com/xxponline/messy-monster-ai-editor/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type RenameAssetReq struct {
	AssetId string `json:"assetId" binding:"required"`
	NewName string `json:"newName" binding:"required"`
}

type RemoveAssetReq struct {
	AssetId string `json:"assetId" binding:"required"`
}

func CreateAssetAPI(context *gin.Context) {
	var req CreateAssetReq
	var initialContent string

	{
		//Basic Request Checking Pass
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
	}

	err := db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error

		//Duplicated Asset Name Checking Pass
		{
			var count int64
			db.GormDatabase.Model(&common.AssetDetailInfo{}).Where("assetSetId = ? AND assetName = ?", req.AssetSetId, req.AssetName).Count(&count)
			if count > 0 {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DuplicatedAssetName,
					"errMessage": common.DuplicatedAssetName.GetMsgFormat(req.AssetName),
				})
				return errors.New(common.DuplicatedAssetName.GetMsgFormat(req.AssetName))
			}
		}
		//Creating Pass
		newAssetId := uuid.New().String()
		{
			newAssetItem := common.AssetDetailInfo{
				AssetId:      newAssetId,
				AssetSetId:   req.AssetSetId,
				AssetName:    req.AssetName,
				AssetType:    req.AssetType,
				AssetVersion: uuid.New().String(),
				AssetContent: initialContent,
			}

			err = db.GormDatabase.Create(&newAssetItem).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}
		//Querying Pass
		var assetItems []common.AssetSummaryInfoItem
		{
			err = db.GormDatabase.Find(&assetItems, "assetSetId = ?", req.AssetSetId).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//All Done
		context.JSON(http.StatusOK, gin.H{
			"errCode":           common.Success,
			"errMessage":        "",
			"assetSummaryInfos": assetItems,
			"newAssetId":        newAssetId,
		})
		return nil
	})

	if err != nil {
		zap.S().Error(err)
	}
}

func RenameAssetAPI(ctx *gin.Context) {

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

	err = db.GormDatabase.Find(&assetItems, "assetSetId = ?", req.AssetSetId).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.DataBaseError,
			"errMessage": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":           common.Success,
		"errMessage":        "",
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

	err = db.GormDatabase.Find(&assetItems, "assetSetId IN ?", req.AssetSetIds).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.DataBaseError,
			"errMessage": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":           common.Success,
		"errMessage":        "",
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

	var assetDetail common.AssetDetailInfo
	err = db.GormDatabase.First(&assetDetail, "id = ?", req.AssetId).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.DataBaseError,
			"errMessage": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":       common.Success,
		"errMessage":    "",
		"assetDocument": assetDetail,
	})
}
