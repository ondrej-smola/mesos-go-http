package agent

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http"
)

type (
	Client struct {
		c mesos.Client
	}
)

func New(c mesos.Client) *Client {
	return &Client{c: c}
}

func (a *Client) GetHealth(ctx context.Context) (*Response_GetHealth, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_HEALTH}, ctx)
	return resp.GetHealth, err
}

func (a *Client) GetFlags(ctx context.Context) (*Response_GetFlags, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_FLAGS}, ctx)
	return resp.GetFlags, err
}

func (a *Client) GetVersion(ctx context.Context) (*Response_GetVersion, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_VERSION}, ctx)
	return resp.GetVersion, err
}

func (a *Client) GetMetrics(req *Call_GetMetrics, ctx context.Context) (*Response_GetMetrics, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_METRICS, GetMetrics: req}, ctx)
	return resp.GetMetrics, err
}

func (a *Client) GetLoggingLevel(ctx context.Context) (*Response_GetLoggingLevel, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_LOGGING_LEVEL}, ctx)
	return resp.GetLoggingLevel, err
}

func (a *Client) SetLoggingLevel(req *Call_SetLoggingLevel, ctx context.Context) error {
	return a.sendWithoutResponse(&Call{Type: Call_SET_LOGGING_LEVEL, SetLoggingLevel: req}, ctx)
}

func (a *Client) ListFiles(req *Call_ListFiles, ctx context.Context) (*Response_ListFiles, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_LIST_FILES, ListFiles: req}, ctx)
	return resp.ListFiles, err
}

func (a *Client) ReadFile(req *Call_ReadFile, ctx context.Context) (*Response_ReadFile, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_READ_FILE, ReadFile: req}, ctx)
	return resp.ReadFile, err
}

func (a *Client) GetState(ctx context.Context) (*Response_GetState, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_STATE}, ctx)
	return resp.GetState, err
}

func (a *Client) GetContainers(ctx context.Context) (*Response_GetContainers, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_CONTAINERS}, ctx)

	// os: Mesos 1.1 returns nil instead of struct with nil slice -  fix it by returning empty struct
	if err == nil {
		conts := resp.GetContainers
		if conts == nil {
			conts = &Response_GetContainers{}
		}
		return conts, nil
	} else {
		return nil, err
	}
}

func (a *Client) GetFrameworks(ctx context.Context) (*Response_GetFrameworks, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_FRAMEWORKS}, ctx)
	return resp.GetFrameworks, err
}

func (a *Client) GetExecutors(ctx context.Context) (*Response_GetExecutors, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_EXECUTORS}, ctx)
	return resp.GetExecutors, err
}

func (a *Client) GetTasks(ctx context.Context) (*Response_GetTasks, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_GET_TASKS}, ctx)
	return resp.GetTasks, err
}

func (a *Client) LaunchNestedContainer(req *Call_LaunchNestedContainer, ctx context.Context) error {
	return a.sendWithoutResponse(&Call{Type: Call_LAUNCH_NESTED_CONTAINER, LaunchNestedContainer: req}, ctx)
}

func (a *Client) WaitNestedContainer(req *Call_WaitNestedContainer, ctx context.Context) (*Response_WaitNestedContainer, error) {
	resp, err := a.sendWithResponse(&Call{Type: Call_WAIT_NESTED_CONTAINER, WaitNestedContainer: req}, ctx)
	return resp.WaitNestedContainer, err
}

func (a *Client) KillNestedContainer(req *Call_KillNestedContainer, ctx context.Context) error {
	return a.sendWithoutResponse(&Call{Type: Call_KILL_NESTED_CONTAINER, KillNestedContainer: req}, ctx)
}

// returns pointer to empty struct in case of error
func (a *Client) sendWithResponse(c *Call, ctx context.Context) (*Response, error) {
	ev := &Response{}
	resp, err := a.c.Do(c, ctx)
	if err != nil {
		return ev, err
	}
	err = resp.Read(ev)
	resp.Close()
	return ev, err
}

func (a *Client) sendWithoutResponse(c *Call, ctx context.Context) error {
	resp, err := a.c.Do(c, ctx)
	if resp != nil {
		resp.Close()
	}
	return err
}
