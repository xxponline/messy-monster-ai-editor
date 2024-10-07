package asset_organization

import (
	"github.com/gin-gonic/gin"
	"github.com/xxponline/messy-monster-ai-editor/asset_content"
	"github.com/xxponline/messy-monster-ai-editor/common"
	"github.com/xxponline/messy-monster-ai-editor/db"
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
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, assetSetMgr := db.ServerDatabase.GetAssetSetManager(true)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}
	defer assetSetMgr.Release()
	var assetSetInfos []common.AssetSetInfoItem
	errCode, errMsg, newAssetSetId := assetSetMgr.CreateAssetSet(req.SolutionId, req.AssetSetName)
	if errCode == common.Success {
		errCode, errMsg, assetSetInfos = assetSetMgr.ListAssetSets(req.SolutionId)
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":       errCode,
		"errMessage":    errMsg,
		"assetSets":     assetSetInfos,
		"newAssetSetId": newAssetSetId,
	})

}

func GetArchivedAssetSetsAPI(context *gin.Context) {
	var req GetAssetSetArchiveReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	if len(req.AssetSetIds) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": "The length of req.AssetSetIds is zero!",
		})
		return
	}

	errCode, errMsg, archivedAssets := doGetArchivedAssetSets(&req)
	if errCode != common.Success {
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
	// ready asset set data
	allArchives := make([]AssetSetArchive, 0, len(req.AssetSetIds))
	{
		errCode, errMsg, setMgr := db.ServerDatabase.GetAssetSetManager(false)
		if errCode != common.Success {
			return errCode, errMsg, nil
		}
		defer setMgr.Release()

		errCode, errMsg, setItems := setMgr.ListAssetSetsBySetIds(req.AssetSetIds)
		if errCode != common.Success {
			return errCode, errMsg, nil
		}

		for _, setItem := range setItems {
			allArchives = append(allArchives, AssetSetArchive{
				AssetSetName: setItem.AssetSetName,
				AssetSetId:   setItem.AssetSetId,
			})
		}
	}

	// process the behaviour trees
	{
		errCode, errMsg, assetMgr := db.ServerDatabase.GetAssetManager(false)
		if errCode != common.Success {
			return errCode, errMsg, nil
		}
		defer assetMgr.Release()

		errCode, errMsg, assetSummaryItems := assetMgr.ListAssets(req.AssetSetIds)

		for i, _ := range allArchives {
			archive := &allArchives[i]
			archivedBtAssets := make([]asset_content.ArchivedBehaviourTree, 0, 8)

			for _, infoItem := range assetSummaryItems {
				if archive.AssetSetId == infoItem.AssetSetId {
					switch infoItem.AssetType {
					case "BehaviourTree":
						errCode, errMsg, archivedBtAsset := asset_content.GetArchivedBehaviourTreeAsset(infoItem.AssetId)
						if errCode != common.Success {
							return errCode, errMsg, nil
						}
						archivedBtAssets = append(archivedBtAssets, *archivedBtAsset)
						break
					default:
						return common.ArchiveAssetsInvalidAssetType, common.ArchiveAssetsInvalidAssetType.GetMsgFormat(infoItem.AssetType), nil
					}
				}
			}
			archive.BehaviourTreeAssets = archivedBtAssets
		}

		return common.Success, "", allArchives
	}
}
