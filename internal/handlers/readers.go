package handlers


func ReadIntoVoid(c chan struct{}) {
	<- c
}