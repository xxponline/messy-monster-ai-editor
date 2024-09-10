package asset_content

import "github.com/gin-gonic/gin"

func InitializeAssetManagement(router *gin.RouterGroup) {
	router.POST("CreateBehaviourTreeNode", CreateBehaviourTreeNodeAPI)
	router.POST("RemoveBehaviourTreeNode", RemoveBehaviourTreeNodeAPI)

}
