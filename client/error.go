package client

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

var (
	UnauthenticatedError = errors.New("Mesos: call not authenticated")
	UnsubscribedError    = errors.New("Mesos: no subscription established")
	VersionError         = errors.New("Mesos: incompatible API version")
	MediaTypeError       = errors.New("Mesos: unsupported media type")
	RateLimitError       = errors.New("Mesos: rate limited")
	UnavailableError     = errors.New("Mesos: server unavailable")
	NotFoundError        = errors.New("Mesos: endpoint not found")
	EmptyResponseErr     = errors.New("Mesos: empty response")
)

type (
	RedirectError struct {
		LeaderHostPort string
	}

	ProtocolError struct {
		err error
	}

	MalformedError struct {
		err error
	}

	// Returning non nil error will fail response processing
	// Should not modify response except when reading response body to return error
	ErrorMapperFunc func(resp *http.Response) error

	//return true if client should retry request
	RetryFunc func(err error) bool
)

func (e RedirectError) Error() string {
	return fmt.Sprintf("Leader changed - is at %v", e.LeaderHostPort)
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("Protocol error, cause %v", e.err)
}

func (e MalformedError) Error() string {
	return fmt.Sprintf("Malformed request: '%v'", e.err)
}

var (
	DefaultErrorMapper = func(resp *http.Response) error {
		code := resp.StatusCode
		switch code {
		case http.StatusOK:
			return nil
		case http.StatusAccepted:
			return nil
		case http.StatusTemporaryRedirect:
			loc := resp.Header.Get("Location")
			if u, err := url.Parse(loc); err != nil {
				return ProtocolError{err: errors.Wrapf(err, "Failed to parse redirect location %v", loc)}
			} else if host, port, err := net.SplitHostPort(u.Host); err != nil {
				return ProtocolError{err: errors.Wrapf(err, "Expected host:port got %v", u.Host)}
			} else {
				return RedirectError{LeaderHostPort: net.JoinHostPort(host, port)}
			}
		case http.StatusBadRequest:
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return MalformedError{err: err}
			} else {
				return MalformedError{err: errors.New(string(body))}
			}
		case http.StatusConflict:
			return VersionError
		case http.StatusForbidden:
			return UnsubscribedError
		case http.StatusUnauthorized:
			return UnauthenticatedError
		case http.StatusNotAcceptable:
			return MediaTypeError
		case http.StatusNotFound:
			return NotFoundError
		case http.StatusServiceUnavailable:
			return UnavailableError
		case http.StatusTooManyRequests:
			return RateLimitError
		default:
			return ProtocolError{err: errors.Errorf("Status code %v", code)}
		}
	}
)

// Returns if errors is Redirect and host:port of redirect location
func IsRedirect(err error) (bool, string) {
	leader := ""

	err = errors.Cause(err)

	if err == nil {
		return false, leader
	}

	l, ok := err.(RedirectError)

	if ok {
		leader = l.LeaderHostPort
	}

	return ok, leader
}
