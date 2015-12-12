package cloud

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var (
	log = logrus.New()
	Log = log
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = &logrus.JSONFormatter{}

	logfile, err := os.Create("debug.log")
	if err != nil {
		log.Fatalln(err)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.Out = logfile

	// Only log the info severity or above.
	log.Level = logrus.InfoLevel

}
