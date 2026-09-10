package bootstrap

import (
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/controller"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/db"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/repository"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/service"
)

type Container struct {
	TransferController *controller.TransferController
}

func InitContainer() *Container {
	// Initialize repositorie
	transferRepo := repository.NewTransferRepository()

	// Initialize service
	transferService := service.NewTransferService(db.DB, transferRepo)

	// Initialize controller
	transferController := controller.NewTransferController(transferService)

	return &Container{
		TransferController: transferController,
	}
}
