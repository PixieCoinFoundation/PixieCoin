package ethereum

func IsAuthenticationNeededErr(err error) bool {
	if _, ok := err.(AuthenticationNeededError); ok {
		return true
	} else {
		return false
	}
}

type AuthenticationNeededError struct {
	Message string
}

func (err AuthenticationNeededError) Error() string {
	return err.Message
}
