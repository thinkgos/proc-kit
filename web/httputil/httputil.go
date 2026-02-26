package httputil

type Response[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId,omitempty"`
	Detail  string `json:"detail,omitempty"`
	Data    T      `json:"data"`
}
