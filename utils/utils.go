package utils

var logger = NewLogger()

// Catch handles Go library errors
func Catch(err error) {
	if err != nil {
		logger.Error(err.Error())
	}
}
