package common

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

//start of the SolutionSummaryInfoItem
//just can be used to quire. modification is forbid

type SolutionSummaryInfoItem struct {
	SolutionId      string `json:"solutionId" binding:"required" gorm:"column:id;primaryKey"`
	SolutionName    string `json:"solutionName" binding:"required" gorm:"column:solutionName"`
	SolutionVersion string `json:"solutionVersion" binding:"required" gorm:"column:solutionVersion"`
}

func (SolutionSummaryInfoItem) TableName() string {
	return "ai_solutions"
}

func (SolutionSummaryInfoItem) BeforeSave(*gorm.DB) error {
	return errors.New("saving SolutionSummaryInfoItem is invalid")
}

func (SolutionSummaryInfoItem) BeforeCreate(*gorm.DB) error {
	return errors.New("creating SolutionSummaryInfoItem is invalid")
}

func (SolutionSummaryInfoItem) BeforeUpdate(*gorm.DB) error {
	return errors.New("updating SolutionSummaryInfoItem is invalid")
}

func (SolutionSummaryInfoItem) BeforeDelete(*gorm.DB) error {
	return errors.New("deleting SolutionSummaryInfoItem is invalid")
}

//end of the SolutionSummaryInfoItem

//start of the SolutionDetailInfo

type SolutionDetailInfo struct {
	SolutionId      string          `json:"solutionId" binding:"required" gorm:"column:id;primaryKey"`
	SolutionName    string          `json:"solutionName" binding:"required" gorm:"column:solutionName"`
	SolutionVersion string          `json:"solutionVersion" binding:"required" gorm:"column:solutionVersion"`
	SolutionMeta    json.RawMessage `json:"solutionMeta" binding:"required" gorm:"column:solutionMeta"`
}

func (SolutionDetailInfo) TableName() string {
	return "ai_solutions"
}

//end of the SolutionDetailInfo

//start of the AssetSetInfoItem

type AssetSetInfoItem struct {
	AssetSetId   string `json:"assetSetId" binding:"required" gorm:"column:id;primaryKey"`
	SolutionId   string `json:"solutionId" binding:"required" gorm:"column:solutionId"`
	AssetSetName string `json:"assetSetName" binding:"required" gorm:"column:assetSetName"`
}

func (AssetSetInfoItem) TableName() string {
	return "ai_asset_sets"
}

//end of the AssetSetInfoItem

//start of the AssetSummaryInfoItem

type AssetSummaryInfoItem struct {
	AssetId      string `json:"assetId" binding:"required" gorm:"column:id;primaryKey"`
	AssetSetId   string `json:"assetSetId" binding:"required" gorm:"column:assetSetId"`
	AssetType    string `json:"assetType" binding:"required" gorm:"column:assetType"`
	AssetName    string `json:"assetName" binding:"required" gorm:"column:assetName"`
	AssetVersion string `json:"assetVersion" binding:"required" gorm:"column:assetVersion"`
}

func (AssetSummaryInfoItem) TableName() string {
	return "ai_asset_documentations"
}

func (AssetSummaryInfoItem) BeforeSave(*gorm.DB) error {
	return errors.New("saving SolutionSummaryInfoItem is invalid")
}

func (AssetSummaryInfoItem) BeforeCreate(*gorm.DB) error {
	return errors.New("creating SolutionSummaryInfoItem is invalid")
}

func (AssetSummaryInfoItem) BeforeUpdate(*gorm.DB) error {
	return errors.New("updating SolutionSummaryInfoItem is invalid")
}

func (AssetSummaryInfoItem) BeforeDelete(*gorm.DB) error {
	return errors.New("deleting SolutionSummaryInfoItem is invalid")
}

//end of the AssetSummaryInfoItem

//start of the AssetDetailInfo

type AssetDetailInfo struct {
	AssetId      string `json:"assetId" binding:"required" gorm:"column:id;primaryKey"`
	AssetSetId   string `json:"assetSetId" binding:"required" gorm:"column:assetSetId"`
	AssetType    string `json:"assetType" binding:"required" gorm:"column:assetType"`
	AssetName    string `json:"assetName" binding:"required" gorm:"column:assetName"`
	AssetVersion string `json:"assetVersion" binding:"required" gorm:"column:assetVersion"`
	AssetContent string `json:"assetContent" binding:"required" gorm:"column:assetContent"`
}

func (AssetDetailInfo) TableName() string {
	return "ai_asset_documentations"
}

//end of the AssetDetailInfo
