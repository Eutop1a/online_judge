package submission

import "online_judge/services"

type ApiGroup struct {
	ApiSubmission
}

var (
	SubmissionService = services.ServiceGroupApp.SubmissionService
)
