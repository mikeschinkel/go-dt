package dt

var userHomeDir DirPath

func EnsureUserHomeDir() (err error) {
	userHomeDir, err = UserHomeDir()
	return err
}

func GetUserHomeDir() DirPath {
	if userHomeDir == "" {
		panic("Must call dt.EnsureUserHomeDir() before dt.GetUserHomeDir() can be called")
	}
	return userHomeDir
}
