package content_modifier

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"messy-monster-ai-editor/common"
	"time"
)

// "bt_root" : BTRootNode, // Not Supported Now
// "bt_selector" : BTSelectorNode,
// "bt_sequence" : BTSequenceNode,
// "bt_simpleParallel" : BTSimpleParallelNode, Not Supported Now
// "bt_task" : BTTaskNode
const (
	Node_Root     = "bt_root"
	Node_Selector = "bt_selector"
	Node_Sequence = "bt_sequence"
	Node_Task     = "bt_task"
)

type XYPosition struct {
	X float32 `json:"x" binding:"required"`
	Y float32 `json:"y" binding:"required"`
}

type LogicBtNode struct {
	NodeId   string          `json:"id" binding:"required"`
	ParentId string          `json:"parentId" binding:"required"`
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
	initializedNodes := []LogicBtNode{{uuid.New().String(), "", XYPosition{100, 100}, Node_Root, 0, nil}}
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
	if nodeType == Node_Selector || nodeType == Node_Sequence || nodeType == Node_Task {
		newNode := LogicBtNode{uuid.New().String(), "", toPosition, nodeType, -1, nil}
		doc.Nodes = append(doc.Nodes, newNode)
		return common.Success, "", &newNode
	}
	return common.BtInvalidNodeType, common.BtInvalidNodeType.GetMsgFormat(nodeType), nil
}

func BehaviourTreeRemoveNode(nodeIds []string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {

	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(nodeIds))
	reserveNodes := make([]LogicBtNode, 0, len(doc.Nodes))
	for _, existNode := range doc.Nodes {
		if slices.Contains(nodeIds, existNode.NodeId) {
			// Need To Removed
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{&existNode, nil})
			if existNode.NodeType == Node_Root {
				return common.BtIllegalRemoveRoot, common.BtIllegalRemoveRoot.GetMsg(), nil
			}

		} else {
			reserveNodes = append(reserveNodes, existNode)
		}
	}
	doc.Nodes = reserveNodes

	return common.Success, "", diffInfos
}

func BehaviourTreeConnectNode(parentId string, childId string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {
	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, 1) // just only one diff when connecting node
	pIdx := -1
	cIdx := -1
	for i, _ := range doc.Nodes {
		if pIdx == -1 && doc.Nodes[i].NodeId == parentId {
			pIdx = i
		}
		if cIdx == -1 && doc.Nodes[i].NodeId == childId {
			cIdx = i
		}

		if pIdx > 0 && cIdx > 0 {
			break
		}
	}

	if pIdx >= 0 && cIdx >= 0 {
		// Check Root Always Not Child
		if doc.Nodes[cIdx].NodeType == Node_Root {
			return common.BtConnectInvalidRootForChild, common.BtConnectInvalidRootForChild.GetMsgFormat(childId), nil
		}
		// Check Task Always Not Parent
		if doc.Nodes[pIdx].NodeType == Node_Task {
			return common.BtConnectInvalidTaskForParent, common.BtConnectInvalidTaskForParent.GetMsgFormat(parentId), nil
		}

		//TODO Front End Has Cycle Check, Consider Do Check In BackEnd

		if doc.Nodes[cIdx].ParentId != parentId { //If The Client Already Connect To The Request Parent, Skip
			preModifiedNode := doc.Nodes[cIdx]
			doc.Nodes[cIdx].ParentId = parentId
			postModifiedNode := doc.Nodes[cIdx]
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{&preModifiedNode, &postModifiedNode})
		}
		return common.Success, "", diffInfos
	} else if pIdx < 0 {
		return common.BtConnectInvalidParent, common.BtConnectInvalidParent.GetMsgFormat(parentId), nil
	} else {
		return common.BtConnectInvalidChild, common.BtConnectInvalidChild.GetMsgFormat(childId), nil
	}

}

func BehaviourTreeDisconnectNode(childIds []string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {
	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(childIds))

	for i, _ := range doc.Nodes {
		if slices.Contains(childIds, doc.Nodes[i].NodeId) {
			if doc.Nodes[i].ParentId == "" {
				return common.BtInvalidDisconnectNodeWithoutParent, common.BtInvalidDisconnectNodeWithoutParent.GetMsgFormat(doc.Nodes[i].NodeId), nil
			}
			//Logic Disconnect
			preModifiedNode := doc.Nodes[i]
			doc.Nodes[i].ParentId = ""
			postModifiedNode := doc.Nodes[i]
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{&preModifiedNode, &postModifiedNode})
		}
	}

	return common.Success, "", diffInfos

}

func BehaviourTreeDisconnectNodeByParentId(parentId string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {
	if parentId == "" {
		return common.BtConnectInvalidParent, common.BtConnectInvalidParent.GetMsgFormat("[Empty]"), nil
	}

	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, 4)
	for i, _ := range doc.Nodes {
		if doc.Nodes[i].ParentId == parentId {
			//Logic Disconnect
			preModifiedNode := doc.Nodes[i]
			doc.Nodes[i].ParentId = ""
			postModifiedNode := doc.Nodes[i]
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{&preModifiedNode, &postModifiedNode})
		}
	}
	return common.Success, "", diffInfos
}
