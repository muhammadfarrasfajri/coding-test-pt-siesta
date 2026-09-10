package route

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/controller"
)

func SetupRouter(r *gin.Engine, transferController *controller.TransferController) {
	r.POST("api/v1/transfer", transferController.HandleTransfer)
}
