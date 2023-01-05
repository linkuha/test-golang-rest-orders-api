package errs

const (
	Other              int = iota // Unclassified error. This value is not printed in the error message.
	InvalidOperation              // Invalid operation for this type of item.
	InvalidArgument               // Invalid argument
	MalformedRequest              // Malformed request body (decode problem).
	IO                            // External I/O error such as network failure.
	Logic                         // Logical error.
	Exist                         // Item already exists.
	NotExist                      // Item does not exist.
	APIAuthorization              // API authorization method related error.
	UserCredentials               // Authentication error (incorrect password, token).
	NotPermitted                  // Has no permissions.
	Private                       // Information withheld.
	Internal                      // Internal error or inconsistency.
	BrokenLink                    // Link target does not exist.
	Database                      // Error from database.
	DatabaseConnection            // Connection to database error.
	RemoteConnection              // Connection to remote service error.
	Validation                    // Input validation error.
	Unanticipated                 // Unanticipated error.
)
