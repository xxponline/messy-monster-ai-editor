package asset_organization

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"messy-monster-ai-editor/common"
	"messy-monster-ai-editor/db"
	"net/http"
)

type CreateSolutionReq struct {
	SolutionName string `json:"solutionName" binding:"required"`
}

func ListSolutionsAPI(context *gin.Context) {

	var errCode common.ErrorCode
	var errMsg string
	var solutionInfos []common.SolutionInfoItem

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
	context.BindJSON(&req)

	//res
	var errCode common.ErrorCode
	var errMsg string
	var solutionInfos []common.SolutionInfoItem

	errCode, errMsg, solutionMgr := db.ServerDatabase.GetSolutionManager(true)
	if errCode == common.Success {
		defer solutionMgr.Release()
		errCode, errMsg = solutionMgr.CreateNewSolution(req.SolutionName)
		if errCode == common.Success {
			errCode, errMsg, solutionInfos = solutionMgr.ListSolutions()
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"errCode":    errCode,
		"errMessage": errMsg,
		"solutions":  solutionInfos,
	})
}
