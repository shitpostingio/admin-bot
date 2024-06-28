package limiter

// AuthorizeUrgentAction authorizes an urgent action
func AuthorizeUrgentAction() {
	urgent <- true
}

// AuthorizeAction authorizes a normal action
func AuthorizeAction() {
	normal <- true
}
