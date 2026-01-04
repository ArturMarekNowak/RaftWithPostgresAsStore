package responses

// This struct is in internal/http/requests but is used only once in postgres_accessor.go in Apply method which might
// look weird. However this struct is a representation of data coming in HTTP request and the Apply method consumes
// *raft.Log as an input. And there is no way to encapsulate it without breaking the interface contract
type DataResponse struct {
	Data string `json:"data"`
}

func NewDataResponse(data string) DataResponse {
	return DataResponse{
		Data: data,
	}
}
