package errorutil

type notFoundError interface {
	NotFound() bool
}

func IsNotFoundError(err error) bool {
	nf, ok := err.(notFoundError)
	return ok && nf.NotFound()
}

type duplicatedError interface {
	Duplicated() bool
}

func IsDuplicatedError(err error) bool {
	d, ok := err.(duplicatedError)
	return ok && d.Duplicated()
}
