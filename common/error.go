package common

import "fmt"

type ErrorCode int

const (
	Success ErrorCode = 0

	//Server or Database Error

	ServerError      ErrorCode = 10000
	DataBaseError    ErrorCode = 10001
	RequestBindError ErrorCode = 10010

	//Asset Organization Error

	DuplicatedSolutionName ErrorCode = 20001
	InvalidSolution        ErrorCode = 20002
	InvalidSolutionVersion ErrorCode = 20003

	DuplicatedAssetSetName ErrorCode = 20010
	InvalidAssetSet        ErrorCode = 20011
	DuplicatedAssetName    ErrorCode = 20021
	InvalidAssetType       ErrorCode = 20022

	//Common Content

	InvalidAssetVersion  ErrorCode = 30001
	DeserializationError ErrorCode = 30010
	SerializationError   ErrorCode = 30011

	//Behaviour Tree Content

	BtInvalidNodeType               ErrorCode = 31001
	BtIllegalRemoveRoot             ErrorCode = 31010
	BtInvalidNodeMovementParameters ErrorCode = 31020

	BtConnectInvalidParent               ErrorCode = 31030
	BtConnectInvalidChild                ErrorCode = 31031
	BtConnectInvalidRootForChild         ErrorCode = 31032
	BtConnectInvalidTaskForParent        ErrorCode = 31033
	BtInvalidDisconnectNodeWithoutParent ErrorCode = 31034
)

var errorMsg = map[ErrorCode]string{
	InvalidSolution:        "Invalid Solution : SolutionId %s ",
	DuplicatedSolutionName: "Duplicated Solution Name %s ",
	InvalidSolutionVersion: "Invalid Solution Version For Modification Exist Version: %s Request Version: %s",

	DuplicatedAssetSetName: "Duplicated AssetSet Name %s",
	InvalidAssetSet:        "Invalid AssetSet : AssetSet Id %s ",
	DuplicatedAssetName:    "Duplicated Asset Name: %s",
	InvalidAssetType:       "Invalid Asset Type: %s",

	InvalidAssetVersion:  "Invalid Asset Version For Modification Exist Version: %s Request Version: %s",
	DeserializationError: "Deserialization Error",
	SerializationError:   "Serialization Error",

	BtInvalidNodeType:               "Invalid Behaviour Tree Node Type: %s",
	BtIllegalRemoveRoot:             "Remove The Root Node In Behaviour Tree Is Illegal",
	BtInvalidNodeMovementParameters: "Invalid Movement Params Length(NodeIds) %d != Length(toPositions) %d",

	BtConnectInvalidParent:               "Invalid Parent Id: %s For Connect/Disconnect",
	BtConnectInvalidChild:                "Invalid Child Id: %s For Connect/Disconnect",
	BtConnectInvalidRootForChild:         "Invalid Parent Id: %s Root Always Not Child",
	BtConnectInvalidTaskForParent:        "Invalid Child Id: %s Task Always Not Parent",
	BtInvalidDisconnectNodeWithoutParent: "Invalid Child Id: %s, Disconnect Node Without Parent",
}

func (errCode ErrorCode) GetMsg() string {
	return errorMsg[errCode]
}

func (errCode ErrorCode) GetMsgFormat(params ...any) string {
	return fmt.Sprintf(errorMsg[errCode], params...)
}
