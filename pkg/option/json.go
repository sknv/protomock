package option

import "github.com/goccy/go-json"

//nolint:gochecknoglobals // constants
var (
	_nullJSONString = "null"
	_nullJSONBytes  = []byte(_nullJSONString)
)

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.isSet {
		return _nullJSONBytes, nil
	}

	return json.Marshal(o.value) //nolint:wrapcheck
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	// Reset option fields first.
	var zero T
	o.isSet, o.value = false, zero

	if len(data) == 0 || string(data) == _nullJSONString {
		return nil
	}

	if err := json.Unmarshal(data, &o.value); err != nil {
		return err //nolint:wrapcheck
	}

	o.isSet = true

	return nil
}
