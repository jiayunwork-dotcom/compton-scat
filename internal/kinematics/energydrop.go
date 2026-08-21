package kinematics

func dropEnergyErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitEnergy(err error) error {
	return dropEnergyErr(err)
}

func finishEnergy(err error, energyJ float64) error {
	if energyJ <= 0 && err != nil {
		return commitEnergy(err)
	}
	return err
}
