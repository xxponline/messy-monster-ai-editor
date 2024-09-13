package common

import "encoding/json"

type SolutionSummaryInfoItem struct {
	SolutionId      string `json:"solutionId" binding:"required"`
	SolutionName    string `json:"solutionName" binding:"required"`
	SolutionVersion string `json:"solutionVersion" binding:"required"`
}

type SolutionDetailInfo struct {
	SolutionSummaryInfoItem
	SolutionMeta json.RawMessage `json:"solutionMeta" binding:"required"`
}

type AssetSetInfoItem struct {
	AssetSetId   string `json:"assetSetId" binding:"required"`
	SolutionId   string `json:"solutionId" binding:"required"`
	AssetSetName string `json:"assetSetName" binding:"required"`
}

type AssetSummaryInfoItem struct {
	AssetId      string `json:"assetId" binding:"required"`
	AssetSetId   string `json:"assetSetId" binding:"required"`
	AssetType    string `json:"assetType" binding:"required"`
	AssetName    string `json:"assetName" binding:"required"`
	AssetVersion string `json:"assetVersion" binding:"required"`
}

type AssetDetailInfo struct {
	AssetId      string `json:"assetId" binding:"required"`
	AssetSetId   string `json:"assetSetId" binding:"required"`
	AssetType    string `json:"assetType" binding:"required"`
	AssetName    string `json:"assetName" binding:"required"`
	AssetContent string `json:"assetContent" binding:"required"`
	AssetVersion string `json:"assetVersion" binding:"required"`
}
