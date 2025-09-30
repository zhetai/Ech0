package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

// Webfinger 处理 Webfinger 请求
func (h *FediverseHandler) Webfinger(ctx *gin.Context) {
	resource := ctx.Query("resource")
	if resource == "" {
		ctx.JSON(http.StatusBadRequest, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   "Missing resource parameter",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// 这里只处理 user 类型的 resource，格式为 acct:username@domain
	var username string
	_, err := fmt.Sscanf(resource, "acct:%s", &username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   "Invalid resource format",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// 调用服务层获取 Actor 信息
	webfingerRes, err := h.service.Webfinger(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   err.Error(),
			Status:  http.StatusNotFound,
		})
		return
	}

	// 设置 Content-Type 为 application/activity+json
	ctx.Header("Content-Type", "application/activity+json")

	// 返回 Actor 信息
	ctx.JSON(http.StatusOK, webfingerRes)
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
			Error:   err.Error(),
			Status:  http.StatusNotFound,
		})
		return
	}

	// 设置 Content-Type 为 application/activity+json
	ctx.Header("Content-Type", "application/activity+json")

	// 返回 Actor 信息
	ctx.JSON(http.StatusOK, actor)
}

// PostInbox 处理接收到的 ActivityPub 消息
func (h *FediverseHandler) PostInbox(ctx *gin.Context) {
	// 从 URL 参数中获取用户名
	username := ctx.Param("username")

	var activity model.Activity
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   "Invalid JSON",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err := h.service.HandleInbox(username, &activity); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// GetOutbox 获取 Outbox 消息
func (h *FediverseHandler) GetOutbox(ctx *gin.Context) {
	// 从 URL 参数中获取用户名
	username := ctx.Param("username")

	var pageStr, pageSizeStr string
	pageStr = ctx.Query("page")
	pageSizeStr = ctx.Query("pageSize")

	// 检查是否带有查询参数
	if pageStr == "" {
		// 不带查询参数，返回 Outbox 元信息
		outbox, err := h.service.BuildOutbox(username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
				Context: "https://www.w3.org/ns/activitystreams",
				Type:    "Error",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}

		ctx.JSON(http.StatusOK, outbox)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	// 校验 page 和 pageSize 的合理范围
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 50 {
		pageSize = 50
	}

	// 打印分页参数
	fmt.Printf("Pagination params - page: %d, pageSize: %d\n", page, pageSize)

	// 调用服务层获取 Outbox 信息
	outbox, err := h.service.HandleOutboxPage(ctx, username, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// 设置 Content-Type 为 application/activity+json
	ctx.Header("Content-Type", "application/activity+json")

	// 返回 Outbox 信息
	ctx.JSON(http.StatusOK, outbox)
}

// GetFollowers 获取粉丝列表
func (h *FediverseHandler) GetFollowers(ctx *gin.Context) {
	// 从 URL 参数中获取用户名
	username := ctx.Param("username")

	// 调用服务层获取粉丝列表
	followers, err := h.service.GetFollowers(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ActivityPubError{
			Context: "https://www.w3.org/ns/activitystreams",
			Type:    "Error",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// 设置 Content-Type 为 application/activity+json
	ctx.Header("Content-Type", "application/activity+json")

	// 返回粉丝列表
	ctx.JSON(http.StatusOK, followers)
}