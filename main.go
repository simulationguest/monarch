package main

func main() {
	config := getConfig()
	selector := NewSelector(config)
	selector.ServeString()
}
