package evaluation

type ApiGroup struct {
	ApiEvaluation
}

func (a *ApiGroup) GetEvaluationApiGroup() ApiEvaluation {
	return a.ApiEvaluation
}
