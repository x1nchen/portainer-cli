package err

// General errors.
const (
	ErrUnauthorized           = Error("Unauthorized")
	ErrResourceAccessDenied   = Error("Access denied to resource")
	ErrAuthorizationRequired  = Error("Authorization required for this operation")
	ErrObjectNotFound         = Error("Object not found inside the database")
	ErrMissingSecurityContext = Error("Unable to find security details in request context")
)

// Error represents an application error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
