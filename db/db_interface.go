package db

import (
	"encoding/json"
	"github.com/xxponline/messy-monster-ai-editor/common"
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
	ListSolutions() (common.ErrorCode, string, []common.SolutionSummaryInfoItem)
	CreateNewSolution(solutionName string) (errCode common.ErrorCode, errMsg string, newSolutionId string)
	ReadSolutionDetail(solutionId string) (errCode common.ErrorCode, errMsg string, solutionInfo *common.SolutionDetailInfo)
	SubmitSolutionMeta(solutionId string, solutionMeta json.RawMessage) (errCode common.ErrorCode, errMsg string, solutionInfo *common.SolutionDetailInfo)
}

type IAssetSetManager interface {
	ListAssetSets(solutionId string) (common.ErrorCode, string, []common.AssetSetInfoItem)
	ListAssetSetsBySetIds(assetSetIds []string) (common.ErrorCode, string, []common.AssetSetInfoItem)
	CreateAssetSet(solutionId string, assetSetName string) (errCode common.ErrorCode, errMsg string, newAssetSetId string)
	Release()
}

type IAssetManager interface {
	ListAssets(assetSetIds []string) (common.ErrorCode, string, []common.AssetSummaryInfoItem)
	CreateAsset(assetSetId string, assetType string, assetName string, assetInitContent string) (errCode common.ErrorCode, errMsg string, createdAssetId string)
	Release()
}

type IAssetDocument interface {
	ReadAsset() (common.ErrorCode, string, *common.AssetDetailInfo)
	UpdateContent(content string) (common.ErrorCode, string, string) //errCode errMsg newVersion
	Release()
}
