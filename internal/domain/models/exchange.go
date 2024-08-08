package models

type Exchange[O any, I any] struct {
	Out O
	In  I
	Err error
}
