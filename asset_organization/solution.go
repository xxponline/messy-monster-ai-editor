package asset_organization

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
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

	errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(false)
	if errCode == common.Success {
		defer solutionMgr.Release()
		errCode, errMsg, solutionInfos = solutionMgr.ListSolutions()
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"solutions":  solutionInfos,
	})
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

	//res
	var errCode common.ErrorCode
	var errMsg string
	var solutionInfos []common.SolutionSummaryInfoItem

	errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(true)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
			"solutions":  solutionInfos,
		})
		return
	}
	defer solutionMgr.Release()
	errCode, errMsg, newSolutionId := solutionMgr.CreateNewSolution(req.SolutionName)
	if errCode != common.Success {
		context.JSON(http.StatusOK, gin.H{
			"errCode":    errCode,
			"errMessage": errMsg,
			"solutions":  solutionInfos,
		})
		return
	}
	errCode, errMsg, solutionInfos = solutionMgr.ListSolutions()

	context.JSON(http.StatusOK, gin.H{
		"errCode":       errCode,
		"errMessage":    errMsg,
		"solutions":     solutionInfos,
		"newSolutionId": newSolutionId,
	})
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

	errCode, errMsg, solutionDetailInfo := doSubmitSolutionMeta(req)

	context.JSON(http.StatusOK, gin.H{
		"errCode":        errCode,
		"errMessage":     errMsg,
		"solutionDetail": solutionDetailInfo,
	})
}

func doSubmitSolutionMeta(req SubmitSolutionMetaReq) (errCode common.ErrorCode, errMsg string, solutionInfo *common.SolutionDetailInfo) {
	errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(true)
	if errCode != common.Success {
		return errCode, errMsg, nil
	}
	defer solutionMgr.Release()

	errCode, errMsg, existDetailInfo := solutionMgr.ReadSolutionDetail(req.SolutionId)
	if errCode != common.Success {
		return errCode, errMsg, nil
	}
	if existDetailInfo.SolutionVersion != req.CurrentVersion {
		return common.InvalidSolutionVersion, common.InvalidSolutionVersion.GetMsgFormat(existDetailInfo.SolutionVersion, req.CurrentVersion), nil
	}

	errCode, errMsg, updatedDetailInfo := solutionMgr.SubmitSolutionMeta(req.SolutionId, req.SolutionMeta)
	return errCode, errMsg, updatedDetailInfo
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

	var errCode common.ErrorCode
	var errMsg string
	var solutionInfo *common.SolutionDetailInfo

	errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(false)
	if errCode == common.Success {
		defer solutionMgr.Release()
		errCode, errMsg, solutionInfo = solutionMgr.ReadSolutionDetail(req.SolutionId)
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":        errCode,
		"errMessage":     errMsg,
		"solutionDetail": solutionInfo,
	})
}
