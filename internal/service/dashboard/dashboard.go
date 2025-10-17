package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	model "github.com/lin-snow/ech0/internal/model/metric"
	"github.com/lin-snow/ech0/internal/monitor"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type DashboardService struct {
	monitor *monitor.Monitor
	commonService    commonService.CommonServiceInterface
}

func NewDashboardService(monitor *monitor.Monitor, commonService commonService.CommonServiceInterface) DashboardServiceInterface {
	return &DashboardService{
		monitor:       monitor,
		commonService: commonService,
	}
}

func (dashboardService *DashboardService) GetMetrics() (model.Metrics, error) {
	return dashboardService.monitor.GetMetrics(), nil
}

func (dashboardService *DashboardService) WSSubsribeMetrics(ctx *gin.Context, userId uint) error {
	// 鉴权
	user, err := dashboardService.commonService.CommonGetUserByUserId(userId)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return nil
	}

	// 接收连接
	conn, err := websocket.Accept(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "Internal Error")


	// 定时推送
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("client disconnected")
			return nil
		case <-ticker.C:
			metrics := dashboardService.monitor.GetMetrics()

			data, _ := json.Marshal(metrics)
			err = conn.Write(ctx, websocket.MessageText, data)
			if err != nil {
				log.Println("write ws error:", err)
				return err
			}
		}
	}
}