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
	ModifiedNodeId   string       `json:"modifiedNodeId" binding:"required"`
	PreModifiedNode  *LogicBtNode `json:"preModifiedNode" binding:"required"`
	PostModifiedNode *LogicBtNode `json:"postModifiedNode" binding:"required"`
}

// in this function, we just use the postModifiedNode of new diff info to replace the same info in exist array
// we are not care the diff info is valid or not (for instance the post modified node is totally same as the previous one)
func mergeOrAppendNodeDiffInfo(nodeDiffInfos []BehaviourTreeNodeDiffInfo, newDiffInfos ...BehaviourTreeNodeDiffInfo) []BehaviourTreeNodeDiffInfo {
	for _, newDiffInfo := range newDiffInfos {
		existNodeIdx := slices.IndexFunc(nodeDiffInfos, func(info BehaviourTreeNodeDiffInfo) bool {
			return info.ModifiedNodeId == newDiffInfo.ModifiedNodeId
		})
		if existNodeIdx > -1 {
			nodeDiffInfos[existNodeIdx].PostModifiedNode = newDiffInfo.PostModifiedNode
		} else {
			nodeDiffInfos = append(nodeDiffInfos, newDiffInfo)
		}
	}
	return nodeDiffInfos
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

	parentIdsForMovedNode := make([]string, 0, len(movementInfos)) //For Reorder
	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(doc.Nodes))

	for nIdx, _ := range doc.Nodes {
		pIdx := slices.IndexFunc(movementInfos, func(item BehaviourTreeNodeMovementItem) bool {
			return item.NodeId == doc.Nodes[nIdx].NodeId
		})
		if pIdx > -1 {
			preModifiedNode := doc.Nodes[nIdx]
			modifyingNode := &doc.Nodes[nIdx]

			modifyingNode.Position = movementInfos[pIdx].ToPosition
			if modifyingNode.ParentId != "" {
				if slices.Index(parentIdsForMovedNode, modifyingNode.ParentId) < 0 {
					parentIdsForMovedNode = append(parentIdsForMovedNode, modifyingNode.ParentId)
				}
			}

			postModifiedNode := *modifyingNode
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{PreModifiedNode: &preModifiedNode, PostModifiedNode: &postModifiedNode})
		}
	}

	//Reorder
	for _, parentId := range parentIdsForMovedNode {
		//Reorder
		diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, reorderBehaviourTreeNodesByParentId(doc, parentId)...)
	}

	return common.Success, "", diffInfos
}

func BehaviourTreeCreateNode(nodeType string, toPosition XYPosition, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {
	//"bt_root" : BTRootNode, // Not Supported Now
	//"bt_selector" : BTSelectorNode,
	//"bt_sequence" : BTSequenceNode,
	//"bt_simpleParallel" : BTSimpleParallelNode, Not Supported Now
	//"bt_task" : BTTaskNode
	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, 1)
	if nodeType == Node_Selector || nodeType == Node_Sequence || nodeType == Node_Task {
		newNode := LogicBtNode{uuid.New().String(), "", toPosition, nodeType, -1, nil}
		doc.Nodes = append(doc.Nodes, newNode)
		diffInfos = []BehaviourTreeNodeDiffInfo{{newNode.NodeId, nil, &newNode}}
		return common.Success, "", diffInfos
	}
	return common.BtInvalidNodeType, common.BtInvalidNodeType.GetMsgFormat(nodeType), nil
}

func BehaviourTreeRemoveNode(nodeIds []string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {

	disconnectedParentIds := make([]string, 0, len(nodeIds)) //For Reorder

	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(nodeIds))
	reserveNodes := make([]LogicBtNode, 0, len(doc.Nodes))
	removedNodeIds := make([]string, 0, len(nodeIds)) // For Disconnect Children
	for _, existNode := range doc.Nodes {
		if slices.Contains(nodeIds, existNode.NodeId) { // Need To Removed
			if existNode.NodeType == Node_Root {
				return common.BtIllegalRemoveRoot, common.BtIllegalRemoveRoot.GetMsg(), nil
			}

			// For Reorder
			if existNode.ParentId != "" {
				if slices.Index(disconnectedParentIds, existNode.ParentId) < 0 {
					disconnectedParentIds = append(disconnectedParentIds, existNode.ParentId)
				}
			}
			removedNodeIds = append(removedNodeIds, existNode.NodeId)
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{existNode.NodeId, &existNode, nil})
		} else {
			reserveNodes = append(reserveNodes, existNode)
		}
	}
	doc.Nodes = reserveNodes

	//Reorder
	for _, parentId := range disconnectedParentIds {
		//Reorder
		diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, reorderBehaviourTreeNodesByParentId(doc, parentId)...)
	}

	//Disconnect Children After Removed Node Because It's Able To Avoid Disconnect Some Node Which Is Removed, It Will Reduce Some Calculate Of DiffInfos Merge
	for _, removedNodeId := range removedNodeIds {
		errCode, errMsg, diffInfosForDisconnect := BehaviourTreeDisconnectNodeByParentId(removedNodeId, doc)
		if errCode != common.Success {
			return errCode, errMsg, nil
		}
		diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, diffInfosForDisconnect...)
	}

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
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{preModifiedNode.NodeId, &preModifiedNode, &postModifiedNode})
		}

		//Reorder
		diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, reorderBehaviourTreeNodesByParentId(doc, parentId)...)

		return common.Success, "", diffInfos
	} else if pIdx < 0 {
		return common.BtConnectInvalidParent, common.BtConnectInvalidParent.GetMsgFormat(parentId), nil
	} else {
		return common.BtConnectInvalidChild, common.BtConnectInvalidChild.GetMsgFormat(childId), nil
	}

}

func BehaviourTreeDisconnectNode(childIds []string, doc *BehaviourTreeDocumentation) (common.ErrorCode, string, []BehaviourTreeNodeDiffInfo) {
	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(childIds))
	disconnectedParentIds := make([]string, 0, len(childIds))
	//Disconnect
	for i, _ := range doc.Nodes {
		if slices.Contains(childIds, doc.Nodes[i].NodeId) {
			if doc.Nodes[i].ParentId == "" {
				return common.BtInvalidDisconnectNodeWithoutParent, common.BtInvalidDisconnectNodeWithoutParent.GetMsgFormat(doc.Nodes[i].NodeId), nil
			}
			//Logic Disconnect
			preModifiedNode := doc.Nodes[i]
			modifyingNode := &doc.Nodes[i]

			if slices.Index(disconnectedParentIds, modifyingNode.ParentId) < 0 {
				disconnectedParentIds = append(disconnectedParentIds, modifyingNode.ParentId)
			}
			modifyingNode.ParentId = ""
			modifyingNode.Order = -1

			postModifiedNode := *modifyingNode
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{preModifiedNode.NodeId, &preModifiedNode, &postModifiedNode})
		}
	}

	//Reorder
	for _, parentId := range disconnectedParentIds {
		//Reorder
		diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, reorderBehaviourTreeNodesByParentId(doc, parentId)...)
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
			modifyingNode := &doc.Nodes[i]

			modifyingNode.ParentId = ""
			modifyingNode.Order = -1
			
			postModifiedNode := *modifyingNode
			diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{preModifiedNode.NodeId, &preModifiedNode, &postModifiedNode})
		}
	}

	//Reorder
	diffInfos = mergeOrAppendNodeDiffInfo(diffInfos, reorderBehaviourTreeNodesByParentId(doc, parentId)...)

	return common.Success, "", diffInfos
}

func reorderBehaviourTreeNodesByParentId(doc *BehaviourTreeDocumentation, parentId string) []BehaviourTreeNodeDiffInfo {
	reorderingNodes := make([]*LogicBtNode, 0, 8)
	for i, _ := range doc.Nodes {
		n := &doc.Nodes[i]
		if n.ParentId == parentId {
			reorderingNodes = append(reorderingNodes, n)
		}
	}

	diffInfos := make([]BehaviourTreeNodeDiffInfo, 0, len(reorderingNodes)/2)
	if len(reorderingNodes) > 0 {
		slices.SortFunc(reorderingNodes, func(a, b *LogicBtNode) int {
			return (int)(a.Position.X - b.Position.X)
		})

		for j, _ := range reorderingNodes {
			if reorderingNodes[j].Order != j {
				preModifiedNode := *reorderingNodes[j]
				reorderingNodes[j].Order = j
				postModifiedNode := *reorderingNodes[j]

				diffInfos = append(diffInfos, BehaviourTreeNodeDiffInfo{preModifiedNode.NodeId, &preModifiedNode, &postModifiedNode})
			}
		}
	}
	return diffInfos
}
