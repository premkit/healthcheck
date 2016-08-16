package healthcheck

type healthcheckUnmarshal struct {
	Name   string            `json:"name"`
	Checks []*checkUnmarshal `json:"checks"`
}

type checkUnmarshal struct {
	Type string                 `json:"type"`
	Spec map[string]interface{} `json:"spec"`
}
