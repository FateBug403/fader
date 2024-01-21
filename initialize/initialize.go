package initialize

func RunInitialize() error {
	var err error
	err = InitViper()
	if err != nil {
		return  err
	}
	return err
}
