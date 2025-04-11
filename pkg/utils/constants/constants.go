package constants

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"os"
)

const (
	IN_PROGRESS          domain.Status           = "IN_PROGRESS"
	COMPLETED            domain.Status           = "COMPLETED"
	REJECTED             domain.Status           = "REJECTED"
	CANCELED             domain.Status           = "CANCELED"
	MAX_FILE_UPLOAD_SIZE                         = int64(3 * 1024 * 1024)
	CLOCK_IN             domain.AttendenceStatus = "CLOCK_IN"
	CLOCK_OUT            domain.AttendenceStatus = "CLOCK_OUT"
)

var (
	PORT             = os.Getenv("PORT")
	APPLICATION_NAME = os.Getenv("APP_NAME")
)
