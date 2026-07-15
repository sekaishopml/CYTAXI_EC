package testing

func Suite(name string, tests ...func()) map[string]func() {
	suite := make(map[string]func(), len(tests))
	for i, t := range tests {
		suite[name] = tests[i]
	}
	return suite
}
