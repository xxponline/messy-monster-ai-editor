package asset_content

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"messy-monster-ai-editor/asset_content/content_modifier"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
	"net/http"
)

type BaseBehaviourTreeModificationReq struct {
	AssetId        string `json:"assetId" binding:"required"`
	CurrentVersion string `json:"currentVersion" binding:"required"`
}

type AssetModifier interface {
	GetAssetID() string
	GetCurrentVersion() string
}

func (req *BaseBehaviourTreeModificationReq) GetAssetID() string {
	return req.AssetId
}

func (req *BaseBehaviourTreeModificationReq) GetCurrentVersion() string {
	return req.CurrentVersion
}

type CreateBehaviourTreeNodeReq struct {
	BaseBehaviourTreeModificationReq
	Position        content_modifier.XYPosition `json:"position" binding:"required"`
	NodeType        string                      `json:"nodeType" binding:"required"`
	InitialSettings json.RawMessage             `json:"initialSettings" binding:"omitempty"`
}

type MoveBehaviourTreeNodeReq struct {
	BaseBehaviourTreeModificationReq
	MovementItems []content_modifier.BehaviourTreeNodeMovementItem `json:"movements" binding:"required"`
}

type RemoveBehaviourTreeNodeReq struct {
	BaseBehaviourTreeModificationReq
	NodeIds []string `json:"nodeIds" binding:"required"`
}

type ConnectBehaviourTreeNodeReq struct {
	BaseBehaviourTreeModificationReq
	ParentNodeId string `json:"parentNodeId" binding:"required"`
	ChildNodeId  string `json:"childNodeId" binding:"required"`
}

type DisconnectBehaviourTreeNodeReq struct {
	BaseBehaviourTreeModificationReq
	ChildNodeIds []string `json:"childNodeIds" binding:"required"`
}

type UpdateBehaviourTreeNodeSettingsReq struct {
	BaseBehaviourTreeModificationReq
	NodeId       string          `json:"nodeId" binding:"required"`
	NodeSettings json.RawMessage `json:"settings" binding:"required"`
}

type GetDetailInfoAboutBehaviourTreeNode struct {
	AssetId string `json:"assetId" binding:"required"`
	NodeId  string `json:"nodeId" binding:"required"`
}

type BehaviourTreeNodeModification struct {
	DiffNodesInfos []content_modifier.BehaviourTreeNodeDiffInfo `json:"diffNodesInfos" binding:"required"`
	PrevVersion    string                                       `json:"prevVersion" binding:"required"`
	NewVersion     string                                       `json:"newVersion" binding:"required"`
}

func CreateBehaviourTreeNodeAPI(context *gin.Context) {
	var req CreateBehaviourTreeNodeReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *CreateBehaviourTreeNodeReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeCreateNode(req.NodeType, req.Position, req.InitialSettings, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func MoveBehaviourTreeNodeAPI(context *gin.Context) {
	var req MoveBehaviourTreeNodeReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *MoveBehaviourTreeNodeReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeMoveNode(req.MovementItems, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func ConnectBehaviourTreeNodeAPI(context *gin.Context) {
	var req ConnectBehaviourTreeNodeReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *ConnectBehaviourTreeNodeReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeConnectNode(req.ParentNodeId, req.ChildNodeId, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func DisconnectBehaviourTreeNodeAPI(context *gin.Context) {
	var req DisconnectBehaviourTreeNodeReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *DisconnectBehaviourTreeNodeReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeDisconnectNode(req.ChildNodeIds, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})

}

func GetDetailInfoAboutBehaviourTreeNodeAPI(context *gin.Context) {
	var req GetDetailInfoAboutBehaviourTreeNode
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, logicNode := DoGetBehaviourTreeNode(req.AssetId, req.NodeId)
	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"nodeInfo":   logicNode,
	})
}

func DoGetBehaviourTreeNode(assetId string, nodeId string) (common.ErrorCode, string, *content_modifier.LogicBtNode) {
	//Db btDoc
	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(assetId, false)
	if errCode != common.Success {
		return errCode, errMsg, nil
	}
	defer assetDoc.Release()

	//detail
	errCode, errMsg, assetDetail := assetDoc.ReadAsset()
	if errCode != common.Success {
		return errCode, errMsg, nil
	}

	//Deserialization
	var btDoc content_modifier.BehaviourTreeDocumentation
	err := json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
	if err != nil {
		return common.DeserializationError, common.DeserializationError.GetMsg(), nil
	}

	for Idx := range btDoc.Nodes {
		if btDoc.Nodes[Idx].NodeId == nodeId {
			// Get It
			gotNode := btDoc.Nodes[Idx]
			return common.Success, "", &gotNode
		}
	}
	return common.BtGetNodeInvalidNodeId, common.BtGetNodeInvalidNodeId.GetMsgFormat(nodeId), nil
}

func UpdateBehaviourTreeNodeSettingsAPI(context *gin.Context) {
	var req UpdateBehaviourTreeNodeSettingsReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *UpdateBehaviourTreeNodeSettingsReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeUpdateNodeSettings(req.NodeId, req.NodeSettings, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func RemoveBehaviourTreeNodeAPI(context *gin.Context) {
	var req RemoveBehaviourTreeNodeReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	errCode, errMsg, modificationInfo := passBehaviourTreeDocumentModification(&req, func(req *RemoveBehaviourTreeNodeReq, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo) {
		return content_modifier.BehaviourTreeRemoveNode(req.NodeIds, btDoc)
	})

	context.JSON(http.StatusOK, gin.H{
		"errCode":          errCode,
		"errMessage":       errMsg,
		"modificationInfo": modificationInfo,
	})
}

func passBehaviourTreeDocumentModification[T AssetModifier](req T, behaviourTreeModify func(req T, btDoc *content_modifier.BehaviourTreeDocumentation) (common.ErrorCode, string, []content_modifier.BehaviourTreeNodeDiffInfo)) (common.ErrorCode, string, *BehaviourTreeNodeModification) {
	//Db btDoc
	errCode, errMsg, assetDoc := db.ServerDatabase.GetAssetDocument(req.GetAssetID(), true)
	if errCode != common.Success {
		return errCode, errMsg, nil
	}
	defer assetDoc.Release()

	//detail
	errCode, errMsg, assetDetail := assetDoc.ReadAsset()
	if errCode != common.Success {
		return errCode, errMsg, nil
	}

	//Version Check
	if assetDetail.AssetVersion != req.GetCurrentVersion() {
		return common.InvalidAssetVersion, common.InvalidAssetVersion.GetMsgFormat(assetDetail.AssetVersion, req.GetCurrentVersion()), nil
	}

	//Deserialization
	var btDoc content_modifier.BehaviourTreeDocumentation
	err := json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
	if err != nil {
		return common.DeserializationError, common.DeserializationError.GetMsg(), nil
	}

	// Real Modified Logic Pass
	errCode, errMsg, diffInfos := behaviourTreeModify(req, &btDoc)
	if errCode != common.Success {
		return errCode, errMsg, nil
	}

	newVersion := req.GetCurrentVersion()
	if len(diffInfos) > 0 { // just need real write data when there are some diffInfos
		//Serialization
		modifiedContent, err := json.Marshal(btDoc)
		if err != nil {
			return common.SerializationError, common.SerializationError.GetMsg(), nil
		}

		//DB Update
		errCode, errMsg, newVersion = assetDoc.UpdateContent(string(modifiedContent))
		if errCode != common.Success {
			return errCode, errMsg, nil
		}
	}

	//Calculate Modification Info
	modificationInfo := BehaviourTreeNodeModification{
		diffInfos,
		req.GetCurrentVersion(),
		newVersion}

	return common.Success, "", &modificationInfo
}
