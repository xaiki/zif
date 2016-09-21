package zif

type AddressResolutionError struct {
	address string
}

func (a AddressResolutionError) Error() string {
	return "Failed to resolve address, address may not exist or is not reachable"
}
