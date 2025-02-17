package k8s

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	Name     string
}

func (e *NotFoundError) Error() string {
	return e.Resource + " " + e.Name + " not found"
}
