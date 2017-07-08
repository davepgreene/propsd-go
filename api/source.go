package api

type Source interface {
	Get()
}

type PollingSource interface {
	Poll()
}
