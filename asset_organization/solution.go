package asset_organization

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xxponline/messy-monster-ai-editor/common"
	"github.com/xxponline/messy-monster-ai-editor/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type CreateSolutionReq struct {
	SolutionName string `json:"solutionName" binding:"required"`
}

type SubmitSolutionMetaReq struct {
	SolutionId     string          `json:"solutionId" binding:"required"`
	CurrentVersion string          `json:"currentVersion" binding:"required"`
	SolutionMeta   json.RawMessage `json:"solutionMeta" binding:"required"`
}

type GetSolutionDetailReq struct {
	SolutionId string `json:"solutionId" binding:"required"`
}

func ListSolutionsAPI(context *gin.Context) {

	var errCode common.ErrorCode
	var errMsg string
	var solutionInfos []common.SolutionSummaryInfoItem

	err := db.GormDatabase.Find(&solutionInfos).Error
	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
			"solutions":  solutionInfos,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.DataBaseError,
			"errMessage": err.Error(),
		})
	}

	//errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(false)
	//if errCode == common.Success {
	//	defer solutionMgr.Release()
	//	errCode, errMsg, solutionInfos = solutionMgr.ListSolutions()
	//}

}

func CreateSolutionAPI(context *gin.Context) {
	//req
	var req CreateSolutionReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	err = db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error
		// Duplicated Name Checking Pass
		{
			var count int64
			tx.Model(&common.SolutionSummaryInfoItem{}).Where("solutionName = ?", req.SolutionName).Count(&count)
			if count > 0 {
				errInfo := common.DuplicatedSolutionName.GetMsgFormat(req.SolutionName)
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DuplicatedSolutionName,
					"errMessage": errInfo,
				})
				return errors.New(errInfo)
			}
		}

		//Create New Solution Pass
		newSolutionId := uuid.New().String()
		{
			newSolutionItem := common.SolutionDetailInfo{
				SolutionId:      newSolutionId,
				SolutionName:    req.SolutionName,
				SolutionVersion: uuid.New().String(),
				SolutionMeta:    json.RawMessage("{}"),
			}

			err = tx.Create(&newSolutionItem).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//Query And Response Pass
		var solutionInfos []common.SolutionSummaryInfoItem
		{
			err := tx.Find(&solutionInfos).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//All Pass Ok
		context.JSON(http.StatusOK, gin.H{
			"errCode":       common.Success,
			"errMessage":    "",
			"solutions":     solutionInfos,
			"newSolutionId": newSolutionId,
		})

		return nil
	})

	if err != nil {
		zap.S().Error(err)
	}
}

func SubmitSolutionMetaAPI(context *gin.Context) {
	var req SubmitSolutionMetaReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	err = db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error
		// Query exist solution item pass
		var existSolutionItem common.SolutionDetailInfo
		{
			err = tx.Find(&existSolutionItem, "id = ?", req.SolutionId).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.InvalidSolution,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//Version checking pass
		{
			if existSolutionItem.SolutionVersion != req.CurrentVersion {
				err = errors.New(common.InvalidSolutionVersion.GetMsgFormat(existSolutionItem.SolutionVersion, req.CurrentVersion))
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.InvalidSolutionVersion,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//Update pass
		{
			existSolutionItem.SolutionMeta = req.SolutionMeta
			err = db.GormDatabase.Save(&existSolutionItem).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.DataBaseError,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		// All Done
		{
			context.JSON(http.StatusOK, gin.H{
				"errCode":        common.Success,
				"errMessage":     "",
				"solutionDetail": existSolutionItem,
			})
		}
		return nil
	})

	if err != nil {
		zap.S().Error(err)
	}
}

func GetSolutionDetailAPI(context *gin.Context) {
	var req GetSolutionDetailReq
	err := context.BindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    common.RequestBindError,
			"errMessage": err.Error(),
		})
		return
	}

	err = db.GormDatabase.Transaction(func(tx *gorm.DB) error {
		var err error
		// Query exist solution item pass
		var existSolutionItem common.SolutionDetailInfo
		{
			err = tx.Find(&existSolutionItem, "id = ?", req.SolutionId).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"errCode":    common.InvalidSolution,
					"errMessage": err.Error(),
				})
				return err
			}
		}

		//All Done
		{
			context.JSON(http.StatusOK, gin.H{
				"errCode":        common.Success,
				"errMessage":     "",
				"solutionDetail": existSolutionItem,
			})
		}
		return nil
	})

	if err != nil {
		zap.S().Error(err)
	}
}
