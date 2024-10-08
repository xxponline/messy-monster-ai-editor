package asset_organization

import "github.com/gin-gonic/gin"

func InitializeAssetManagement(router *gin.RouterGroup) {
	router.GET("ListSolutions", ListSolutionsAPI)
	router.POST("CreateSolution", CreateSolutionAPI)
	router.POST("GetSolutionDetail", GetSolutionDetailAPI)
	router.POST("SubmitSolutionMeta", SubmitSolutionMetaAPI)

	router.POST("ListAssetSets", ListAssetSetsAPI)
	router.POST("CreateAssetSet", CreateAssetSetAPI)
	router.POST("GetArchivedAssetSets", GetArchivedAssetSetsAPI)

	router.POST("CreateAsset", CreateAssetAPI)
	router.POST("ListAssets", ListAssetsAPI)
	router.POST("ListAssetsByMultipleAssetSets", ListAssetsByMultipleAssetSetsAPI)
	router.POST("ReadAsset", ReadAssetAPI)
}
