package reportController

import (
	"time"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/services/reportService"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type reportListElement struct {
	ID           uint   `json:"id"`
	ReporterName string `json:"reporter_name"`
	SellerName   string `json:"seller_name"`
	BookName     string `json:"book_name"`
	Status       uint   `json:"status"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
}

// GetReportList 管理员获取举报列表
func GetReportList(c *gin.Context) {
	list, err := reportService.GetReportList()
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	reportList := make([]reportListElement, 0)
	for _, report := range list {
		reportList = append(reportList, reportListElement{
			ID:           report.ID,
			ReporterName: report.ReporterName,
			SellerName:   report.SellerName,
			BookName:     report.BookName,
			Status:       report.Status,
			Title:        report.Title,
			CreatedAt:    report.CreatedAt.Format(time.DateTime),
		})
	}

	response.JsonSuccessResp(c, gin.H{
		"report_list": reportList,
	})
}
