package leader

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"github.com/ondrej-smola/mesos-go-http/log"
	"github.com/pkg/errors"
	"net/url"
	"sync"
)

var NoAvailableLeaderFoundErr = errors.New("Mesos: request failed on all leaders")

type (
	Opt func(*LeaderClient)

	// Client for handling requests that must be send to leading master.
	LeaderClient struct {
		clientOpts     []client.Opt
		clientProvider client.Provider
		maxRedirects   int
		masters        []string
		log            log.Logger

		sync.RWMutex
		// mutex for
		leader   mesos.Client
		endpoint string
	}
)

func WithClientProvider(p client.Provider) Opt {
	return func(l *LeaderClient) {
		l.clientProvider = p
	}
}

func WithClientOpts(opts ...client.Opt) Opt {
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
		log:          log.NewNopLogger(),
		masters:      endpoints,
		maxRedirects: 5,
	}

	for _, o := range opts {
		o(l)
	}

	if l.clientProvider == nil {
		l.clientProvider = client.NewProvider(l.clientOpts...)
	}

	return l
}

// Sends Message to current leader (following redirects) and returning response.
// Current leader is reused for subsequent requests.
func (c *LeaderClient) Do(msg codec.Message, ctx context.Context, opts ...mesos.RequestOpt) (mesos.Response, error) {
	c.RLock()
	leader := c.leader
	endpoint := c.endpoint
	c.RUnlock()

	// shared path
	if leader != nil {
		if resp, err := leader.Do(msg, ctx, opts...); err == nil {
			return resp, err
		} else {
			// request failed -> find new leader
			c.log.Log("event", "call", "err", err)
		}
	}

	// slow path
	c.Lock()
	defer c.Unlock()
	c.leader = nil
	c.endpoint = ""

	for _, m := range c.masters {
		endpoint = m

		for i := 0; i < c.maxRedirects; i++ {
			newClient := c.clientProvider(endpoint)
			if resp, err := newClient.Do(msg, ctx, opts...); err != nil {
				if ok, newLeader := client.IsRedirect(err); ok {
					leaderUrl, err := url.Parse(endpoint)
					if err != nil {
						return nil, errors.Wrap(err, "URL parse")
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
				c.leader = newClient
				c.endpoint = endpoint
				return resp, nil
			}
		}
	}

	return nil, NoAvailableLeaderFoundErr
}
