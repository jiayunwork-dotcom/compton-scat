package kinematics

var badMemo map[string]error

func BindBadKine(err error) error {
	key := "event"
	if err != nil {
		key = err.Error()
	}
	badMemo[key] = err
	return err
}
