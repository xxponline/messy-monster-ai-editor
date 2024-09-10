package content_modifier

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"messy-monster-ai-editor/common"
	"time"
)

type XYPosition struct {
	X float32 `json:"x" binding:"required"`
	Y float32 `json:"y" binding:"required"`
}

type LogicBtNode struct {
	NodeId   string          `json:"id" binding:"required"`
	Position XYPosition      `json:"position" binding:"required"`
	NodeType string          `json:"type" binding:"required"`
	Order    int             `json:"order" binding:"required"`
	Data     json.RawMessage `json:"data" binding:"required"`
}

type LogicBtConnection struct {
	ConnectionId string `json:"id" binding:"required"`
	Source       string `json:"source" binding:"required"`
	Target       string `json:"target" binding:"required"`
}

type LogicBtDescriptor struct {
	DescriptorId string `json:"id" binding:"required"`
	AttachTo     string `json:"attachTo" binding:"required"`
	Order        int    `json:"order" binding:"required"`
}

type LogicBtService struct {
	ServiceId string `json:"id" binding:"required"`
	AttachTo  string `json:"attachTo" binding:"required"`
	Order     int    `json:"order" binding:"required"`
}

type BehaviourTreeDocumentation struct {
	ModifyTimeStamp int64
	Nodes           []LogicBtNode
	Connections     []LogicBtConnection
	Descriptors     []LogicBtDescriptor
	Services        []LogicBtService
}

type BehaviourTreeNodeMovementItem struct {
	NodeId     string     `json:"nodeId" binding:"required"`
	ToPosition XYPosition `json:"toPosition" binding:"required"`
}

type BehaviourTreeNodeDiffInfo struct {
	PreModifiedNode  *LogicBtNode `json:"preModifiedNode" binding:"required"`
	PostModifiedNode *LogicBtNode `json:"postModifiedNode" binding:"required"`
}

func BehaviourTreeCreateEmptyContent() (errCode common.ErrorCode, errMsg string, content string) {
	initializedNodes := []LogicBtNode{{uuid.New().String(), XYPosition{100, 100}, "bt_root", 0, nil}}
	BehaviourTreeDocumentation := BehaviourTreeDocumentation{
		ModifyTimeStamp: time.Now().UTC().Unix(),
		Nodes:           initializedNodes,
		Connections:     []LogicBtConnection{},
		Descriptors:     []LogicBtDescriptor{},
		Services:        []LogicBtService{},
	}

	b, err := json.Marshal(BehaviourTreeDocumentation)
	if err != nil {
		return common.SerializationError, common.SerializationError.GetMsg(), ""
	}
	return common.Success, "", string(b)
}

func BehaviourTreeMoveNode(movementInfos []BehaviourTreeNodeMovementItem, doc *BehaviourTreeDocumentation) (errCode common.ErrorCode, errMsg string, updateNodeDiffs []BehaviourTreeNodeDiffInfo) {

	resultDiff := make([]BehaviourTreeNodeDiffInfo, 0, len(doc.Nodes))

	for nIdx, _ := range doc.Nodes {
		pIdx := slices.IndexFunc(movementInfos, func(item BehaviourTreeNodeMovementItem) bool {
			return item.NodeId == doc.Nodes[nIdx].NodeId
		})
		if pIdx > -1 {
			preUpdateNode := doc.Nodes[nIdx]
			doc.Nodes[nIdx].Position = movementInfos[pIdx].ToPosition
			afterUpdateNode := doc.Nodes[nIdx]
			resultDiff = append(resultDiff, BehaviourTreeNodeDiffInfo{PreModifiedNode: &preUpdateNode, PostModifiedNode: &afterUpdateNode})
		}
	}
	return common.Success, "", resultDiff
}

func BehaviourTreeCreateNode(nodeType string, toPosition XYPosition, doc *BehaviourTreeDocumentation) (errCode common.ErrorCode, errMsg string, createdNode *LogicBtNode) {
	//"bt_root" : BTRootNode, // Not Supported Now
	//"bt_selector" : BTSelectorNode,
	//"bt_sequence" : BTSequenceNode,
	//"bt_simpleParallel" : BTSimpleParallelNode, Not Supported Now
	//"bt_task" : BTTaskNode
	if nodeType == "bt_selector" || nodeType == "bt_sequence" || nodeType == "bt_task" {
		newNode := LogicBtNode{uuid.New().String(), toPosition, nodeType, -1, nil}
		doc.Nodes = append(doc.Nodes, newNode)
		return common.Success, "", &newNode
	}
	return common.BtInvalidNodeType, common.BtInvalidNodeType.GetMsgFormat(nodeType), nil
}

func BehaviourTreeRemoveNode(nodeIds []string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []LogicBtNode) {
	reserveNodes := make([]LogicBtNode, 0, len(doc.Nodes))
	removedNodes := make([]LogicBtNode, 0, len(doc.Nodes))
	for _, existNode := range doc.Nodes {
		if slices.Contains(nodeIds, existNode.NodeId) {
			// Need To Removed
			removedNodes = append(removedNodes, existNode)

			if existNode.NodeType == "bt_root" {
				return common.BtIllegalRemoveRoot, common.BtIllegalRemoveRoot.GetMsg(), nil
			}

		} else {
			reserveNodes = append(reserveNodes, existNode)
		}
	}
	doc.Nodes = reserveNodes

	return common.Success, "", removedNodes
}
