package publish

type Publisher interface {
	Configured() bool
	Publish() error
}
