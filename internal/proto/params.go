package proto

type StartBackupRequest struct {
	IsCompress   bool
	ResourceType string
	ResourceOpts []byte
	StuType      string
	StuOpts      []byte
}

type StartBackupResponse struct {
	JobID int
}

type GetJobStatusRequest struct {
	JobID int
}

type GetJobResultResponse struct {
	Status string
	Result []byte
}
