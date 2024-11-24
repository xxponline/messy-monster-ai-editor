package asset_organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xxponline/messy-monster-ai-editor/asset_content"
	"github.com/xxponline/messy-monster-ai-editor/common"
	"github.com/xxponline/messy-monster-ai-editor/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type ListAssetSetReq struct {
	SolutionId string `json:"solutionId" binding:"required"`
}

type CreateAssetSetReq struct {
	SolutionId   string `json:"solutionId" binding:"required"`
	AssetSetName string `json:"assetSetName" binding:"required"`
}

type GetAssetSetArchiveReq struct {
	AssetSetIds []string `json:"assetSetIds" binding:"required"`
}

type AssetSetArchive struct {
	AssetSetName        string                                `json:"assetSetName" binding:"required"`
	AssetSetId          string                                `json:"assetSetId" binding:"required"`
	BehaviourTreeAssets []asset_content.ArchivedBehaviourTree `json:"behaviourTreeAssets" binding:"required"`
}

func ListAssetSetsAPI(context *gin.Context) {
	var req ListAssetSetReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	var assetSetInfos []common.AssetSetInfoItem

	err = db.GormDatabase.Find(&assetSetInfos).Error

	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.DataBaseError,
			"errMessage": err.Error(),
		})
		zap.S().Error(err)
	} else {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.Success,
			"errMessage": "",
			"assetSets":  assetSetInfos,
		})
	}
}

func CreateAssetSetAPI(context *gin.Context) {
	var req CreateAssetSetReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error
		//Duplicated Name Checking Pass
		{
			var count int64
			tx.Model(&common.AssetSetInfoItem{}).Where("solutionId = ? AND assetSetName = ?", req.SolutionId, req.AssetSetName).Count(&count)
			if count > 0 {
				errMsg := common.DuplicatedAssetSetName.GetMsgFormat(req.AssetSetName)
				errors.New(errMsg)
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DuplicatedSolutionName,
					"errMessage": errMsg,
				})
				return err
			}
		}

		//Create New Asset Set Item Pass
		newAssetSetId := uuid.New().String()
		{
			newAssetSetItem := common.AssetSetInfoItem{
				AssetSetId:   newAssetSetId,
				AssetSetName: req.AssetSetName,
				SolutionId:   req.SolutionId,
			}

			err = tx.Create(&newAssetSetItem).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//Query Pass
		var assetSetInfos []common.AssetSetInfoItem
		{
			err := tx.Find(&assetSetInfos).Error
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
			"errCode":       common.Success,
			"errMessage":    "",
			"assetSets":     assetSetInfos,
			"newAssetSetId": newAssetSetId,
		})

		return nil
	})

	if err != nil {
		zap.S().Error(err)
	}
}

func GetArchivedAssetSetsAPI(context *gin.Context) {
	var req GetAssetSetArchiveReq
	err := context.BindJSON(&req)
	if err != nil {
		zap.S().Warn(err)
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	if len(req.AssetSetIds) == 0 {
		zap.S().Warn("the length of req.AssetSetIds is zero")
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": "The length of req.AssetSetIds is zero!",
		})
		return
	}

	errCode, errMsg, archivedAssets := doGetArchivedAssetSets(&req)
	if errCode != common.Success {
		zap.S().Warn(errMsg)
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":        common.Success,
		"errMessage":     "",
		"archivedAssets": archivedAssets,
	})

}

func doGetArchivedAssetSets(req *GetAssetSetArchiveReq) (common.ErrorCode, string, []AssetSetArchive) {
	var errCode common.ErrorCode
	var errMsg string
	var allArchives []AssetSetArchive

	errMsg = ""
	errCode = common.Success

	err := db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error
		// ready asset set data
		allArchives = make([]AssetSetArchive, 0, len(req.AssetSetIds))
		{
			var assetSets []common.AssetSetInfoItem
			err = tx.Find(&assetSets, "id IN ?", req.AssetSetIds).Error
			if err != nil {
				return err
			}

			for _, setItem := range assetSets {
				allArchives = append(allArchives, AssetSetArchive{
					AssetSetName: setItem.AssetSetName,
					AssetSetId:   setItem.AssetSetId,
				})
			}
		}

		// process the behaviour trees
		{
			var assetItems []common.AssetDetailInfo

			err = tx.Find(&assetItems).Error
			if err != nil {
				return err
			}

			for i, _ := range allArchives {
				archive := &allArchives[i]
				archivedBtAssets := make([]asset_content.ArchivedBehaviourTree, 0, 8)

				for itemIdx, _ := range assetItems {
					if archive.AssetSetId == assetItems[itemIdx].AssetSetId {
						switch assetItems[itemIdx].AssetType {
						case "BehaviourTree":
							var archivedBtAsset *asset_content.ArchivedBehaviourTree
							errCode, errMsg, archivedBtAsset = asset_content.GetArchivedBehaviourTreeAsset(&assetItems[itemIdx])
							if errCode != common.Success {
								return errors.New(errMsg)
							}
							archivedBtAssets = append(archivedBtAssets, *archivedBtAsset)
							break
						default:
							errCode = common.ArchiveAssetsInvalidAssetType
							errMsg = common.ArchiveAssetsInvalidAssetType.GetMsgFormat(&assetItems[itemIdx].AssetType)
							return errors.New(errMsg)
						}
					}
				}
				archive.BehaviourTreeAssets = archivedBtAssets
			}
		}

		return nil
	})

	if err != nil {
		return errCode, errMsg, nil
	}
	return common.Success, "", allArchives
}
