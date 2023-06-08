package services

type ServiceStatus interface {
	Error() error
	OK() bool
}

type Service interface {
	Start() error
	Stop()
	Status() ServiceStatus
}
