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

func (s *DashboardService) WSSubsribeMetrics(w http.ResponseWriter, r *http.Request) error {
	// 鉴权
	// user, err := s.commonService.CommonGetUserByUserId(userId)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte("unauthorized"))
	// 	return err
	// }
	// if !user.IsAdmin {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	w.Write([]byte("permission denied"))
	// 	return nil
	// }

	// WebSocket 升级
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return err
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer conn.Close()

		for {
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

			data, _ := json.Marshal(resp)

			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("write ws error:", err)
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()
	return nil
}
