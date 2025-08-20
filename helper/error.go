package helper

func PanicError(e error) {
	if e != nil {
		panic(e)
	}
}

func PrintAndMoveError(e error) {
	if e != nil {
		print(e)
	}
}
