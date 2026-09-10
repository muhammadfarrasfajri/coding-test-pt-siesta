package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/model"
	errs "github.com/muhammadfarrasfajri/coding-test-pt-siesta/pkg"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/service"
)

type TransferController struct {
	ServiceTransfer service.TransferService
}

func NewTransferController(serviceTransfer service.TransferService) *TransferController {
	return &TransferController{
		ServiceTransfer: serviceTransfer,
	}
}

func (c *TransferController) HandleTransfer(ctx *gin.Context) {

	var req model.TransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Invalid request payload",
			Type:    "InvalidPayload",
		})
		return
	}

	idempotencyKey := ctx.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "X-Idempotency-Key header is required",
			Type:    "MissingHeader",
		})
		return
	}

	traceID := ctx.GetHeader("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.New().String()
	}

	err := c.ServiceTransfer.ExecuteTransfer(ctx.Request.Context(), traceID, idempotencyKey, req.FromAccount, req.ToAccount, req.Amount)

	if err != nil {

		if coreErr, ok := err.(*errs.CoreBankingError); ok {
			// Jika error karena idempotency conflict, anggap sukses
			if coreErr.Code == "ERR_IDEMPOTENT_CONFLICT" {
				ctx.JSON(http.StatusOK, model.APIResponse{
					Error:   true,
					Message: "Transfer already processed",
					Type:    "IdempotentConflict",
					Data: gin.H{
						"trace_id": traceID,
					},
				})
				return
			}

			// Tentukan HTTP Status berdasarkan Severity
			httpStatus := http.StatusInternalServerError
			if coreErr.Severity == errs.SeverityWarn {
				httpStatus = http.StatusBadRequest
			}

			ctx.JSON(httpStatus, model.APIResponse{
				Error:   true,
				Message: coreErr.Message,
				Type:    coreErr.Code,
				Data: gin.H{
					"trace_id":  coreErr.TraceID,
					"retryable": coreErr.IsRetryable,
				},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Error:   true,
			Message: "Internal server error",
			Type:    "InternalServerError",
			Data: gin.H{
				"trace_id": traceID,
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Transfer successful",
		Data: gin.H{
			"trace_id": traceID,
		},
	})
}
