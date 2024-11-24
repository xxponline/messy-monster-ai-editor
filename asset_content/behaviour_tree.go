package asset_content

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xxponline/messy-monster-ai-editor/asset_content/content_modifier"
	"github.com/xxponline/messy-monster-ai-editor/common"
	"github.com/xxponline/messy-monster-ai-editor/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type ArchivedBehaviourTree struct {
	AssetName    string `json:"assetName" binding:"required"`
	AssetId      string `json:"assetId" binding:"required"`
	AssetVersion string `json:"assetVersion" binding:"required"`

	BehaviourTreeNodes       string `json:"behaviourTreeNodes" binding:"required"`
	BehaviourTreeDescriptors string `json:"behaviourTreeDescriptors" binding:"required"`
	BehaviourTreeServices    string `json:"behaviourTreeServices" binding:"required"`
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

	errCode, errMsg, logicNode := doGetBehaviourTreeNode(req.AssetId, req.NodeId)
	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"nodeInfo":   logicNode,
	})
}

func doGetBehaviourTreeNode(assetId string, nodeId string) (common.ErrorCode, string, *content_modifier.LogicBtNode) {

	var assetDetail common.AssetDetailInfo
	err := db.GormDatabase.First(&assetDetail, "id = ?", assetId).Error
	if err != nil {
		return common.DataBaseError, err.Error(), nil
	}

	//Deserialization
	var btDoc content_modifier.BehaviourTreeDocumentation
	err = json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
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
	var errCode = common.Success
	var errMsg = ""
	var modificationInfo *BehaviourTreeNodeModification = nil

	err := db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var assetDetail common.AssetDetailInfo
		//Querying Pass
		{
			err := db.GormDatabase.First(&assetDetail, "id = ?", req.GetAssetID()).Error
			if err != nil {
				errCode = common.DataBaseError
				errMsg = err.Error()
				return err
			}
		}

		//Version Checking Pass
		{
			if assetDetail.AssetVersion != req.GetCurrentVersion() {
				eMsg := common.InvalidAssetVersion.GetMsgFormat(assetDetail.AssetVersion, req.GetCurrentVersion())
				errCode, errMsg = common.InvalidAssetVersion, eMsg
				return errors.New(eMsg)
			}
		}

		//Deserialization Pass
		var btDoc content_modifier.BehaviourTreeDocumentation
		{
			err := json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
			if err != nil {
				eMsg := common.DeserializationError.GetMsg()
				errCode, errMsg = common.DeserializationError, eMsg
				return errors.New(eMsg)
			}
		}

		//Real Modified Logic Pass
		var diffInfos []content_modifier.BehaviourTreeNodeDiffInfo
		{
			errCode, errMsg, diffInfos = behaviourTreeModify(req, &btDoc)
			if errCode != common.Success {
				return errors.New(errMsg)
			}
		}

		//Write Modification
		newVersion := req.GetCurrentVersion()
		if len(diffInfos) > 0 { // just need real write data when there are some diffInfos
			//Serialization
			modifiedContent, err := json.Marshal(btDoc)
			if err != nil {
				errCode = common.SerializationError
				eMsg := errCode.GetMsg()
				errMsg = eMsg
				return errors.New(eMsg)
			}

			//DB Update
			newVersion = uuid.New().String()
			assetDetail.AssetVersion = newVersion
			assetDetail.AssetContent = string(modifiedContent)

			err = tx.Save(assetDetail).Error
			if err != nil {
				errCode = common.DataBaseError
				eMsg := errCode.GetMsg()
				errMsg = eMsg
				return errors.New(eMsg)
			}
		}

		//All Pass
		//Calculate Modification Info
		modificationInfo = &BehaviourTreeNodeModification{
			diffInfos,
			req.GetCurrentVersion(),
			newVersion}

		return nil
	})

	if err != nil {
		zap.S().Error(err)
		return errCode, errMsg, nil
	}
	return common.Success, "", modificationInfo
}

func GetArchivedBehaviourTreeAsset(assetDetail *common.AssetDetailInfo) (common.ErrorCode, string, *ArchivedBehaviourTree) {

	if assetDetail.AssetType != "BehaviourTree" {
		return common.ArchiveAssetsUnexpectAssetType, common.ArchiveAssetsUnexpectAssetType.GetMsgFormat(assetDetail.AssetType, "BehaviourTree"), nil
	}

	var btDoc content_modifier.BehaviourTreeDocumentation
	//Deserialize
	err := json.Unmarshal([]byte(assetDetail.AssetContent), &btDoc)
	if err != nil {
		return common.DeserializationError, err.Error(), nil
	}

	var archivedDoc ArchivedBehaviourTree
	archivedDoc.AssetName = assetDetail.AssetName
	archivedDoc.AssetId = assetDetail.AssetId
	archivedDoc.AssetVersion = assetDetail.AssetVersion

	archivedNodes := make([]map[string]interface{}, 0, len(btDoc.Nodes))
	for _, node := range btDoc.Nodes {
		archivedNode := make(map[string]interface{})
		//fmt.Printf("%v %d \n", node.Settings, len(node.Settings))
		archivedNode["id"] = node.NodeId
		archivedNode["type"] = node.NodeType
		archivedNode["order"] = node.Order
		archivedNode["parentId"] = node.ParentId
		//fmt.Printf("before desirialize %v \n", archivedNode)
		if string(node.Settings) != "null" {
			fmt.Printf("umarshal ... %s \n", string(node.Settings))
			err := json.Unmarshal(node.Settings, &archivedNode)
			if err != nil {
				return common.DeserializationError, err.Error(), nil
			}
		}
		//fmt.Printf("after desirialize %v \n", archivedNode)

		////
		//fmt.Printf("%v \n", archivedNode)

		archivedNodes = append(archivedNodes, archivedNode)
	}

	{
		serializationNodes, err := json.Marshal(&archivedNodes)
		if err != nil {
			return common.SerializationError, err.Error(), nil
		}

		archivedDoc.BehaviourTreeNodes = string(serializationNodes)
	}

	archivedDoc.BehaviourTreeDescriptors = "[]"
	archivedDoc.BehaviourTreeServices = "[]"
	return common.Success, "", &archivedDoc

}
