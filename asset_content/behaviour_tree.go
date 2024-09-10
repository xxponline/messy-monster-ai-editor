package asset_content

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"messy-monster-ai-editor/asset_content/content_modifier"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
	"net/http"
)

type CreateBehaviourTreeNodeReq struct {
	AssetId        string                      `json:"assetId" binding:"required"`
	CurrentVersion string                      `json:"currentVersion" binding:"required"`
	Position       content_modifier.XYPosition `json:"position" binding:"required"`
	NodeType       string                      `json:"nodeType" binding:"required"`
}

type MoveBehaviourTreeNodeReq struct {
	AssetId        string                      `json:"assetId" binding:"required"`
	NodeId         string                      `json:"nodeId" binding:"required"`
	CurrentVersion string                      `json:"currentVersion" binding:"required"`
	Position       content_modifier.XYPosition `json:"position" binding:"required"`
}

type RemoveBehaviourTreeNodeReq struct {
	AssetId        string   `json:"assetId" binding:"required"`
	NodeIds        []string `json:"nodeIds" binding:"required"`
	CurrentVersion string   `json:"currentVersion" binding:"required"`
}

type UpdateBehaviourTreeNodeSettingsReq struct {
	AssetId        string          `json:"assetId" binding:"required"`
	NodeId         string          `json:"nodeId" binding:"required"`
	CurrentVersion string          `json:"currentVersion" binding:"required"`
	NodeSettings   json.RawMessage `json:"data" binding:"required"`
}

type BehaviourTreeNodeDiffInfo struct {
	PreModifiedNode  *content_modifier.LogicBtNode `json:"preModifiedNode" binding:"required"`
	PostModifiedNode *content_modifier.LogicBtNode `json:"postModifiedNode" binding:"required"`
}

type BehaviourTreeNodeModification struct {
	DiffNodesInfos []BehaviourTreeNodeDiffInfo `json:"diffNodesInfos" binding:"required"`
	PrevVersion    string                      `json:"prevVersion" binding:"required"`
	NewVersion     string                      `json:"newVersion" binding:"required"`
}

func CreateBehaviourTreeNodeAPI(context *gin.Context) {
	var req CreateBehaviourTreeNodeReq
	context.BindJSON(&req)

	//Db doc
	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(req.AssetId, true)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}
	defer assetDoc.Release()

	//detail
	errCode, errMsg, assetDetail := assetDoc.ReadAsset()
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Version Check
	if assetDetail.AssetVersion != req.CurrentVersion {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.InvalidAssetVersion,
			"errMessage": common.InvalidAssetVersion.GetMsgFormat(assetDetail.AssetVersion, req.CurrentVersion),
		})
		return
	}

	//Deserialization
	var doc content_modifier.BehaviourTreeDocumentation
	{
		err := json.Unmarshal([]byte(assetDetail.AssetContent), &doc)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"errCode":    errCode,
				"errMessage": errMsg,
			})
			return
		}
	}

	//Real Create
	errCode, errMsg, createdNode := content_modifier.BehaviourTreeCreateNode(req.NodeType, req.Position, &doc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Serialization
	modifiedContent, err := json.Marshal(doc)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//DB Update
	errCode, errMsg, newVersion := assetDoc.UpdateContent(string(modifiedContent))
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Calculate Modification Info
	modificationInfo := BehaviourTreeNodeModification{
		[]BehaviourTreeNodeDiffInfo{{nil, createdNode}},
		req.CurrentVersion,
		newVersion}

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func MoveBehaviourTreeNodeAPI(context *gin.Context) (int, string, *BehaviourTreeNodeModification) {

	return 0, "", nil
}

func DoMoveBehaviourTreeNode(req MoveBehaviourTreeNodeReq) (int, string, *BehaviourTreeNodeModification) {
	return 0, "", nil
}

func LinkBehaviourTreeNode() {

}

func ModifyBehaviourTreeNode() {

}

func RemoveBehaviourTreeNodeAPI(context *gin.Context) {
	var req RemoveBehaviourTreeNodeReq
	context.BindJSON(&req)

	//Db doc
	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(req.AssetId, true)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}
	defer assetDoc.Release()

	//detail
	errCode, errMsg, assetDetail := assetDoc.ReadAsset()
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}
	
	//Version Check
	if assetDetail.AssetVersion != req.CurrentVersion {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.InvalidAssetVersion,
			"errMessage": common.InvalidAssetVersion.GetMsgFormat(assetDetail.AssetVersion, req.CurrentVersion),
		})
		return
	}

	//Deserialization
	var doc content_modifier.BehaviourTreeDocumentation
	{
		err := json.Unmarshal([]byte(assetDetail.AssetContent), &doc)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"errCode":    errCode,
				"errMessage": errMsg,
			})
			return
		}
	}

	//Real Remove
	errCode, errMsg, removedNodes := content_modifier.BehaviourTreeRemoveNode(req.NodeIds, &doc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Serialization
	modifiedContent, err := json.Marshal(doc)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//DB Update
	errCode, errMsg, newVersion := assetDoc.UpdateContent(string(modifiedContent))
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Calculate Modification Info
	diffInfos := make([]BehaviourTreeNodeDiffInfo, len(removedNodes))
	for _, removedNode := range removedNodes {
		diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{&removedNode, nil})
	}
	modificationInfo := BehaviourTreeNodeModification{
		diffInfos,
		req.CurrentVersion,
		newVersion}

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}
