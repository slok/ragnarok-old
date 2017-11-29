package controller

// Controller is the piece that will know how to process the informer "informs". A controller
// needs a way of getting informed about the changes in the resources so it can apply a logic, this
// can get done using an informer. Depending the nature of the informer, it will need to use different
// approaces.
type Controller interface {
	// Run will run the controller and ensure the state of the resources is the desired one.
	Run() error
	// Stop will stop the controller.
	Stop() error
}
