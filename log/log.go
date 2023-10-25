package log

import "go.uber.org/zap"

func NewLog() (*zap.Logger, error) {
	c := zap.NewProductionConfig()

	//c.OutputPaths = []string{"stderr", "./snippetbox.log"}
	return c.Build()
}
