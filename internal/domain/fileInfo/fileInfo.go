package fileInfo

const (
	Available        FileInfoStatus = "Available"
	NotAvailable     FileInfoStatus = "NotAvailable"
	NotAllowedFormat FileInfoStatus = "NotAllowedFormat"
	Loaded           FileInfoStatus = "Loaded"
	NotLoaded        FileInfoStatus = "NotLoaded"
)

type FileInfo struct {
	URL    string
	Status FileInfoStatus
}

type FileInfoStatus string
