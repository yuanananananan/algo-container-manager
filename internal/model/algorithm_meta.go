package model

type AlgorithmMeta struct {
	AlgorithmID   uint64 `json:"algorithmId"`
	AlgorithmCode string `json:"algorithmCode"`
	Name          string `json:"name"`
	AlgorithmType string `json:"algorithmType"`
	Framework     string `json:"framework"`
	RuntimeType   string `json:"runtimeType"`
}
