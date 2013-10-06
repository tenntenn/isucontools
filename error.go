package isucontools

func Must(err error) {
	if err != nil {
		log.Panic(err)
	}
}
