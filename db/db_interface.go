package db

import (
	"messy-monster-ai-editor/common"
)

type IAiEditorDatabase interface {
	Initialize(dataSource string) (common.ErrorCode, string)
	GetSolutionManager(WriteLock bool) (common.ErrorCode, string, ISolutionManager)
	GetAssetSetManager(WriteLock bool) (common.ErrorCode, string, IAssetSetManager)
	GetAssetManager(WriteLock bool) (common.ErrorCode, string, IAssetManager)
	GetAssetDocument(AssetId string, WriteLock bool) (common.ErrorCode, string, IAssetDocument)
}

type ISolutionManager interface {
	Release()
	ListSolutions() (common.ErrorCode, string, []common.SolutionInfoItem)
	CreateNewSolution(solutionName string) (common.ErrorCode, string)
}

type IAssetSetManager interface {
	ListAssetSets(solutionId string) (common.ErrorCode, string, []common.AssetSetInfoItem)
	CreateAssetSet(solutionId string, assetSetName string) (common.ErrorCode, string)
	Release()
}

type IAssetManager interface {
	ListAssets(assetSetId string) (common.ErrorCode, string, []common.AssetSummaryInfoItem)
	CreateAsset(assetSetId string, assetType string, assetName string, assetInitContent string) (common.ErrorCode, string)
	Release()
}

type IAssetDocument interface {
	ReadAsset() (common.ErrorCode, string, *common.AssetDetailInfo)
	UpdateContent(content string) (common.ErrorCode, string, string) //errCode errMsg newVersion
	Release()
}
