package leader

import (
	"context"
	"fmt"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"github.com/ondrej-smola/mesos-go-http/log"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

var NoAvailableLeaderFoundErr = errors.New("Mesos: unable to connect to any leader")

type (
	Opt func(*LeaderClient)

	// Client for handling requests that must be send to leading master.
	LeaderClient struct {
		clientOpts     []client.Opt
		endpointFunc   mesos.EndpointFunc
		clientProvider client.Provider
		maxRedirects   int
		masters        mesos.Masters

		log log.Logger

		sync.RWMutex
		// mutex for
		leader mesos.Client
	}
)

func WithClientProvider(p client.Provider) Opt {
	return func(l *LeaderClient) {
		l.clientProvider = p
	}
}

func WithMasters(hostPorts ...string) Opt {
	masters := mesos.NewMasters(hostPorts...)
	if err := masters.Valid(); err != nil {
		panic(fmt.Sprintf("Invalid mesos masters: %v", err))
	}

	return func(l *LeaderClient) {
		l.masters = masters
	}
}

func WithClientOpts(opts ...client.Opt) Opt {
	return func(l *LeaderClient) {
		l.clientOpts = opts
	}
}

func WithEndpointFunc(f mesos.EndpointFunc) Opt {
	return func(l *LeaderClient) {
		l.endpointFunc = f
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

func New(opts ...Opt) *LeaderClient {
	l := &LeaderClient{
		endpointFunc: mesos.V1ApiEndpointFunc,
		log:          log.NewNopLogger(),
		masters:      mesos.DEFAULT_MASTERS,
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

	c.log.Log("event", "connecting to leader", "to", strings.Join(c.masters, ","))

	for _, m := range c.masters {
		endpoint := c.endpointFunc(m)

		for i := 0; i < c.maxRedirects; i++ {
			newClient := c.clientProvider(endpoint)
			if resp, err := newClient.Do(msg, ctx, opts...); err != nil {
				if ok, newLeader := client.IsRedirect(err); ok {
					to := c.endpointFunc(newLeader)
					c.log.Log("event", "redirected", "from", endpoint, "to", to, "attempt", i+1)
					endpoint = to
					continue
				} else {
					c.log.Log("event", "call", "endpoint", endpoint, "err", err)
					break
				}
			} else {
				c.log.Log("event", "connected_to_leader", "endpoint", endpoint)
				c.leader = newClient
				return resp, nil
			}
		}
	}

	return nil, NoAvailableLeaderFoundErr
}
