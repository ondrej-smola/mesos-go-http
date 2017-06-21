package leader

import (
	"context"
	"net/url"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/lib/client"
	"github.com/ondrej-smola/mesos-go-http/lib/log"
	"github.com/pkg/errors"
)

type (
	Opt func(*LeaderClient)

	// Client for handling requests that must be send to leading master
	LeaderClient struct {
		clientOpts     []client.DefaultClientOpt
		clientProvider client.Provider
		maxRedirects   int
		masters        []string
		log            log.Logger

		sync.RWMutex

		// mutex for
		endpoint    string
		endpointSeq uint64 // sequence number for parallel leader find sync
	}
)

func WithClientProvider(p client.Provider) Opt {
	return func(l *LeaderClient) {
		l.clientProvider = p
	}
}

func WithClientOpts(opts ...client.DefaultClientOpt) Opt {
	return func(l *LeaderClient) {
		l.clientOpts = opts
	}
}

// Maximum number of redirects during single connection attempt to leader.
func WithMaxRedirects(count int) Opt {
	return func(l *LeaderClient) {
		l.maxRedirects = count
	}
}

func WithLogger(l log.Logger) Opt {
	return func(c *LeaderClient) {
		c.log = l
	}
}

func New(endpoints []string, opts ...Opt) *LeaderClient {
	l := &LeaderClient{
		log:            log.NewNopLogger(),
		masters:        endpoints,
		clientProvider: client.NewProvider(),
		maxRedirects:   5,
	}

	for _, o := range opts {
		o(l)
	}

	return l
}

// Send message to current leader.
// Handles leader (re)detection and retries on error
func (c *LeaderClient) Do(msg proto.Message, ctx context.Context, opts ...client.RequestOpt) (client.Response, error) {
	for {
		c.RLock()
		endpoint := c.endpoint
		lastSeq := c.endpointSeq
		c.RUnlock()
		// concurrent path
		if endpoint != "" {
			leader := c.clientProvider.New(endpoint, c.clientOpts...)
			if resp, err := leader.Do(msg, ctx, opts...); err == nil {
				return resp, err
			} else {
				// request failed -> find new leader
				c.log.Log("event", "call", "err", err)
			}
		}

		c.Lock()

		// somebody already found new leader
		if c.endpointSeq > lastSeq {
			c.Unlock()
			continue
		}

		if resp, newEndpoint, err := c.findLeader(msg, ctx, opts...); err != nil {
			c.Unlock()
			return nil, err
		} else {
			c.endpoint = newEndpoint
			c.endpointSeq = lastSeq + 1
			c.Unlock()
			return resp, nil
		}
	}
}

func (c *LeaderClient) findLeader(msg proto.Message, ctx context.Context, opts ...client.RequestOpt) (client.Response, string, error) {
	var lastErr error

	for _, m := range c.masters {
		endpoint := m

		for i := 0; i < c.maxRedirects; i++ {
			select {
			case <-ctx.Done():
				return nil, "", ctx.Err()
			default:
			}

			leader := c.clientProvider.New(endpoint, c.clientOpts...)
			if resp, err := leader.Do(msg, ctx, opts...); err != nil {
				lastErr = err
				if ok, newLeader := client.IsRedirect(err); ok {
					leaderUrl, err := url.Parse(endpoint)
					if err != nil {
						return nil, "", errors.Wrap(err, "URL parse")
					}
					leaderUrl.Host = newLeader
					to := leaderUrl.String()
					c.log.Log(
						"event", "redirected",
						"from", endpoint,
						"to", to,
						"attempt", i+1,
						"debug", true)
					endpoint = to
					continue
				} else {
					c.log.Log("event", "call", "endpoint", endpoint, "err", err)
					break
				}
			} else {
				c.log.Log(
					"event", "connected_to_leader",
					"endpoint", endpoint,
					"debug", true)
				return resp, endpoint, nil
			}
		}
	}

	return nil, "", errors.Errorf("Mesos: request failed on all endpoints, last error: %v", lastErr)
}
