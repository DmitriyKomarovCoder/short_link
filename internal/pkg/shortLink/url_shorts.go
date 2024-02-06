package shortLink

type Usecase interface {
	GetUrl(url string) (string, error)
	CreateLink(url string) (string, error)
}
