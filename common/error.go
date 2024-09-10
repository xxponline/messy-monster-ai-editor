package common

import "fmt"

type ErrorCode int

const (
	Success ErrorCode = 0

	//Server or Database Error

	ServerError   ErrorCode = 10000
	DataBaseError ErrorCode = 10001

	//Asset Organization Error

	DuplicatedSolutionName ErrorCode = 20001
	InvalidSolution        ErrorCode = 20002
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
)

var errorMsg = map[ErrorCode]string{
	InvalidSolution:        "Invalid Solution : SolutionId %s ",
	DuplicatedSolutionName: "Duplicated Solution Name %s ",
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
}

func (errCode ErrorCode) GetMsg() string {
	return errorMsg[errCode]
}

func (errCode ErrorCode) GetMsgFormat(params ...any) string {
	return fmt.Sprintf(errorMsg[errCode], params)
}