package bootstrap

import "github.com/muhammadfarrasfajri/coding-test-pt-siesta/db"

func InitDatabase() {
	db.Connect()
}
