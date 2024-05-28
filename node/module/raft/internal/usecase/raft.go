package usecase

type RaftUsecase interface {
	Get(key any) (string, error)
	Set(key any, val string) (string, error)
	Strlen(key any) (int, error)
	Del(key any) (string, error)
	Append(key any, value string) (string, error)
}
