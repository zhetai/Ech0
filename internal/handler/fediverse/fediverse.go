package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	service "github.com/lin-snow/ech0/internal/service/fediverse"
)

type FediverseHandler struct {
	service service.FediverseServiceInterface
}

func NewFediverseHandler(fediverseService service.FediverseServiceInterface) *FediverseHandler {
	return &FediverseHandler{
		service: fediverseService,
	}
}

// GetActor 获取 Actor 信息
func (h *FediverseHandler) GetActor(ctx *gin.Context) {
	// 从 URL 参数中获取用户名
	username := ctx.Param("username")

	// 调用服务层获取 Actor 信息
	actor, err := h.service.GetActorByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error: err.Error(),
			Status: http.StatusNotFound,
		})
		return
	}

	// 设置 Content-Type 为 application/activity+json
	ctx.Header("Content-Type", "application/activity+json")

	// 返回 Actor 信息
	ctx.JSON(http.StatusOK, actor)
}