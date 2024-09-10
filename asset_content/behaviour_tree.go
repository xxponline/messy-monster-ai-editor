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
	AssetId        string                                           `json:"assetId" binding:"required"`
	MovementItems  []content_modifier.BehaviourTreeNodeMovementItem `json:"movements" binding:"required"`
	CurrentVersion string                                           `json:"currentVersion" binding:"required"`
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

type BehaviourTreeNodeModification struct {
	DiffNodesInfos []content_modifier.BehaviourTreeNodeDiffInfo `json:"diffNodesInfos" binding:"required"`
	PrevVersion    string                                       `json:"prevVersion" binding:"required"`
	NewVersion     string                                       `json:"newVersion" binding:"required"`
}

func CreateBehaviourTreeNodeAPI(context *gin.Context) {
	var req CreateBehaviourTreeNodeReq
	context.BindJSON(&req)

	errCode, errMsg, assetDoc, btDoc := readyBehaviourTreeDocForUpdate(req.AssetId, req.CurrentVersion)
	defer assetDoc.Release()
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Real Create
	errCode, errMsg, createdNode := content_modifier.BehaviourTreeCreateNode(req.NodeType, req.Position, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	errCode, errMsg, newVersion := writeBehaviourTreeDocForUpdate(assetDoc, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Calculate Modification Info
	modificationInfo := BehaviourTreeNodeModification{
		[]content_modifier.BehaviourTreeNodeDiffInfo{{nil, createdNode}},
		req.CurrentVersion,
		newVersion}

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func MoveBehaviourTreeNodeAPI(context *gin.Context) {
	var req MoveBehaviourTreeNodeReq
	context.BindJSON(&req)

	errCode, errMsg, assetDoc, btDoc := readyBehaviourTreeDocForUpdate(req.AssetId, req.CurrentVersion)
	defer assetDoc.Release()
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Real Create
	errCode, errMsg, nodeDiffInfos := content_modifier.BehaviourTreeMoveNode(req.MovementItems, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	// write To DB
	errCode, errMsg, newVersion := writeBehaviourTreeDocForUpdate(assetDoc, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	modificationInfo := BehaviourTreeNodeModification{
		nodeDiffInfos,
		req.CurrentVersion,
		newVersion}

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func LinkBehaviourTreeNode() {

}

func ModifyBehaviourTreeNode() {

}

func RemoveBehaviourTreeNodeAPI(context *gin.Context) {
	var req RemoveBehaviourTreeNodeReq
	context.BindJSON(&req)

	errCode, errMsg, assetDoc, btDoc := readyBehaviourTreeDocForUpdate(req.AssetId, req.CurrentVersion)
	defer assetDoc.Release()
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Real Remove
	errCode, errMsg, diffInfos := content_modifier.BehaviourTreeRemoveNode(req.NodeIds, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	errCode, errMsg, newVersion := writeBehaviourTreeDocForUpdate(assetDoc, btDoc)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
		})
		return
	}

	//Calculate Modification Info
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

func readyBehaviourTreeDocForUpdate(assetId string, requestVersion string) (common.ErrorCode, string, db.IAssetDocument, *content_modifier.BehaviourTreeDocumentation) {
	//Db btDoc
	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(assetId, true)
	if errCode != common.Success {
		return errCode, errMsg, nil, nil
	}

	//detail
	errCode, errMsg, assetDetail := assetDoc.ReadAsset()
	if errCode != common.Success {
		return errCode, errMsg, nil, nil
	}

	//Version Check
	if assetDetail.AssetVersion != requestVersion {
		return common.InvalidAssetVersion, common.InvalidAssetVersion.GetMsgFormat(assetDetail.AssetVersion, requestVersion), nil, nil
	}

	//Deserialization
	var btDoc content_modifier.BehaviourTreeDocumentation
	err := json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
	if err != nil {
		return common.DeserializationError, common.DeserializationError.GetMsg(), nil, nil
	}

	return common.Success, "", assetDoc, &btDoc
}

func writeBehaviourTreeDocForUpdate(assetDoc db.IAssetDocument, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, string) { //result errCode errMsg, newVersion
	//Serialization
	modifiedContent, err := json.Marshal(btDoc)
	if err != nil {
		return common.SerializationError, common.SerializationError.GetMsg(), ""
	}

	//DB Update
	errCode, errMsg, newVersion := assetDoc.UpdateContent(string(modifiedContent))
	if errCode != common.Success {
		return errCode, errMsg, ""
	}
	return common.Success, "", newVersion
}
