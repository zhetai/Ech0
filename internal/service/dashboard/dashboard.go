package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	model "github.com/lin-snow/ech0/internal/model/metric"
	"github.com/lin-snow/ech0/internal/monitor"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type DashboardService struct {
	monitor       *monitor.Monitor
	commonService commonService.CommonServiceInterface
}

func NewDashboardService(
	monitor *monitor.Monitor,
	commonService commonService.CommonServiceInterface,
) DashboardServiceInterface {
	return &DashboardService{
		monitor:       monitor,
		commonService: commonService,
	}
}

func (dashboardService *DashboardService) GetMetrics() (model.Metrics, error) {
	return dashboardService.monitor.GetMetrics(), nil
}

func (s *DashboardService) WSSubsribeMetrics(w http.ResponseWriter, r *http.Request, userId uint) error {
	// 鉴权
	user, err := s.commonService.CommonGetUserByUserId(userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return err
	}
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("permission denied"))
		return nil
	}

	// WebSocket 升级
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			log.Printf("client disconnected")
			return nil
		case <-ticker.C:
			metrics := s.monitor.GetMetrics()

			resp := struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
				Data any    `json:"data"`
			}{
				Code: 1,
				Msg:  "metrics update",
				Data: metrics,
			}

			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("json marshal error:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("write ws error: %v", err)
				return err
			}
		}
	}
}
