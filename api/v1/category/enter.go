package category

import "online_judge/services"

type ApiGroup struct {
	ApiCategory
}

var (
	CategoryService = services.ServiceGroupApp.CategoryService
)
