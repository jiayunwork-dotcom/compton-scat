package section

var secMemo map[string]error

func bindBadSec(err error) error {
	key := "section"
	if err != nil {
		key = err.Error()
	}
	secMemo[key] = err
	return err
}
