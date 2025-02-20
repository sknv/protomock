package option

import "fmt"

type Option[T any] struct {
	value T
	isSet bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		isSet: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{} //nolint:exhaustruct // empty by default
}

func (o Option[T]) IsSome() bool {
	return o.isSet
}

func (o Option[T]) IsNone() bool {
	return !o.isSet
}

func (o Option[T]) Unwrap() T { //nolint:ireturn // generic method
	return o.Expect("Called 'Unwrap' on a 'None' value")
}

func (o Option[T]) Expect(panicMsg string) T { //nolint:ireturn // generic method
	if o.isSet {
		return o.value
	}

	panic(panicMsg)
}

func (o Option[T]) UnwrapOr(other T) T { //nolint:ireturn // generic method
	if o.isSet {
		return o.value
	}

	return other
}

func (o Option[T]) UnwrapOrDefault() T { //nolint:ireturn // generic method
	if o.isSet {
		return o.value
	}

	var def T

	return def
}

func (o Option[T]) UnwrapOrElse(other func() T) T { //nolint:ireturn // generic method
	if o.isSet {
		return o.value
	}

	return other()
}

func (o Option[T]) String() string {
	if o.isSet {
		return fmt.Sprintf("Some(%v)", o.value)
	}

	return "None"
}
