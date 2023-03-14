package abi

func AddType(name, signature string) error {
	var err error
	Default.Types[name], err = ParseType(signature)
	return err
}

func MustAddType(name, signature string) {
	if err := AddType(name, signature); err != nil {
		panic(err)
	}
}
