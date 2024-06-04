package fetcher

type Fetcher interface {
	FetchData() (any, error)
}
