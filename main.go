package main

func main() {
	println(foo())
}

func foo() (bar string) {
	defer func() {
		if r := recover(); r != nil {
			bar = "recovered from panic"
		}
	}()
	defer func() {
		bar = "deferred"
	}()
	bar = "initial value"
	return
}
