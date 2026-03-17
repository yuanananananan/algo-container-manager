package api

import "github.com/gin-gonic/gin"

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	// 1. 创建
	r.POST("/containers/start", handler.StartAlgorithm)

	// 2. 查询（管理视角：数据库）
	r.GET("/containers", handler.ListContainers)

	// 3.  查询（运行视角：K8s）
	r.GET("/containers/runtime", handler.ListRuntimeContainers)
	r.GET("/containers/:name/status", handler.GetContainerStatus)

	// 4. 操作（生命周期管理）
	r.DELETE("/containers/:name", handler.DeleteContainer)
	r.POST("/containers/:name/restart", handler.RestartContainer)
	r.POST("/containers/:name/scale", handler.ScaleContainer)

	return r
}
