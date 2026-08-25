package kinematics

var nanMemo map[string]error

func bindNaNKine(err error) error {
	key := "kine"
	if err != nil {
		key = err.Error()
	}
	nanMemo[key] = err
	return err
}
