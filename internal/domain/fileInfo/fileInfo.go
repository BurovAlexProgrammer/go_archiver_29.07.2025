package fileInfo

const (
	Available        Status = "Available"
	NotAvailable     Status = "NotAvailable"
	NotAllowedFormat Status = "NotAllowedFormat"
	Loaded           Status = "Loaded"
	NotLoaded        Status = "NotLoaded"
)

type FileInfo struct {
	URL    string
	Status Status
}

type Status string
