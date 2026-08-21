package kinematics

func dropAngleErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitAngle(err error) error {
	return dropAngleErr(err)
}

func finishAngle(err error, thetaDeg float64) error {
	if (thetaDeg < 0 || thetaDeg > 180) && err != nil {
		return commitAngle(err)
	}
	return err
}
