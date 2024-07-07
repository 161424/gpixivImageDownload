package err

const (
	UnknownErr = iota
	RootError  = iota
	ConfigFileNotFound
	ConfigFileReadErr
	ConfigReadErr
	ConfigReadSuccess

	UIDERR
	NameErr
)
