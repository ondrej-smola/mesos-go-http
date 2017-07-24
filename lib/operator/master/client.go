package master

import (
	"context"

	"github.com/ondrej-smola/mesos-go-http/lib/client"
)

type (
	Client struct {
		c client.Client
	}

	EventStream struct {
		c client.Client
	}

	EventOrErr struct {
		Event *Event
		Err   error
	}
)

func New(c client.Client) *Client {
	return &Client{c: c}
}

func (a *Client) GetHealth(ctx context.Context) (*Response_GetHealth, error) {
	resp, err := a.call(&Call{Type: Call_GET_HEALTH.Enum()}, ctx)
	return resp.GetHealth, err
}

func (a *Client) GetFlags(ctx context.Context) (*Response_GetFlags, error) {
	resp, err := a.call(&Call{Type: Call_GET_FLAGS.Enum()}, ctx)
	return resp.GetFlags, err
}

func (a *Client) GetVersion(ctx context.Context) (*Response_GetVersion, error) {
	resp, err := a.call(&Call{Type: Call_GET_VERSION.Enum()}, ctx)
	return resp.GetVersion, err
}

func (a *Client) GetMetrics(req *Call_GetMetrics, ctx context.Context) (*Response_GetMetrics, error) {
	resp, err := a.call(&Call{Type: Call_GET_METRICS.Enum(), GetMetrics: req}, ctx)
	return resp.GetMetrics, err
}

func (a *Client) GetLoggingLevel(ctx context.Context) (*Response_GetLoggingLevel, error) {
	resp, err := a.call(&Call{Type: Call_GET_LOGGING_LEVEL.Enum()}, ctx)
	return resp.GetLoggingLevel, err
}

func (a *Client) SetLoggingLevel(req *Call_SetLoggingLevel, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_SET_LOGGING_LEVEL.Enum(), SetLoggingLevel: req}, ctx)
}

func (a *Client) ListFiles(req *Call_ListFiles, ctx context.Context) (*Response_ListFiles, error) {
	resp, err := a.call(&Call{Type: Call_LIST_FILES.Enum(), ListFiles: req}, ctx)
	return resp.ListFiles, err
}

func (a *Client) ReadFile(req *Call_ReadFile, ctx context.Context) (*Response_ReadFile, error) {
	resp, err := a.call(&Call{Type: Call_READ_FILE.Enum(), ReadFile: req}, ctx)
	return resp.ReadFile, err
}

func (a *Client) GetState(ctx context.Context) (*Response_GetState, error) {
	resp, err := a.call(&Call{Type: Call_GET_STATE.Enum()}, ctx)
	return resp.GetState, err
}

func (a *Client) GetAgents(ctx context.Context) (*Response_GetAgents, error) {
	resp, err := a.call(&Call{Type: Call_GET_AGENTS.Enum()}, ctx)
	return resp.GetAgents, err
}

func (a *Client) GetFrameworks(ctx context.Context) (*Response_GetFrameworks, error) {
	resp, err := a.call(&Call{Type: Call_GET_FRAMEWORKS.Enum()}, ctx)
	return resp.GetFrameworks, err
}

func (a *Client) GetExecutors(ctx context.Context) (*Response_GetExecutors, error) {
	resp, err := a.call(&Call{Type: Call_GET_EXECUTORS.Enum()}, ctx)
	return resp.GetExecutors, err
}

func (a *Client) GetTasks(ctx context.Context) (*Response_GetTasks, error) {
	resp, err := a.call(&Call{Type: Call_GET_TASKS.Enum()}, ctx)
	return resp.GetTasks, err
}

func (a *Client) GetRoles(ctx context.Context) (*Response_GetRoles, error) {
	resp, err := a.call(&Call{Type: Call_GET_ROLES.Enum()}, ctx)
	return resp.GetRoles, err
}

func (a *Client) GetWeights(ctx context.Context) (*Response_GetWeights, error) {
	resp, err := a.call(&Call{Type: Call_GET_WEIGHTS.Enum()}, ctx)

	// os: Mesos 1.1.0 returns nil instead of struct with nil slice -  fix it by returning empty struct
	if err == nil {
		weights := resp.GetWeights
		if weights == nil {
			weights = &Response_GetWeights{}
		}
		return weights, nil
	} else {
		return nil, err
	}
}

func (a *Client) UpdateWeights(req *Call_UpdateWeights, ctx context.Context) (*Response_GetWeights, error) {
	resp, err := a.call(&Call{Type: Call_UPDATE_WEIGHTS.Enum(), UpdateWeights: req}, ctx)
	return resp.GetWeights, err
}

func (a *Client) GetMaster(ctx context.Context) (*Response_GetMaster, error) {
	resp, err := a.call(&Call{Type: Call_GET_MASTER.Enum()}, ctx)
	return resp.GetMaster, err
}

func (a *Client) ReserveResources(req *Call_ReserveResources, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_RESERVE_RESOURCES.Enum(), ReserveResources: req}, ctx)
}

func (a *Client) UnreserveResources(req *Call_UnreserveResources, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_UNRESERVE_RESOURCES.Enum(), UnreserveResources: req}, ctx)
}

func (a *Client) CreateVolumes(req *Call_CreateVolumes, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_CREATE_VOLUMES.Enum(), CreateVolumes: req}, ctx)
}

func (a *Client) DestroyVolumes(req *Call_DestroyVolumes, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_DESTROY_VOLUMES.Enum(), DestroyVolumes: req}, ctx)
}

func (a *Client) GetMaintenanceStatus(ctx context.Context) (*Response_GetMaintenanceStatus, error) {
	resp, err := a.call(&Call{Type: Call_GET_MAINTENANCE_STATUS.Enum()}, ctx)
	return resp.GetMaintenanceStatus, err
}

func (a *Client) GetMaintenanceSchedule(ctx context.Context) (*Response_GetMaintenanceSchedule, error) {
	resp, err := a.call(&Call{Type: Call_GET_MAINTENANCE_SCHEDULE.Enum()}, ctx)
	return resp.GetMaintenanceSchedule, err
}

func (a *Client) UpdateMaintenanceSchedule(req *Call_UpdateMaintenanceSchedule, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_UPDATE_MAINTENANCE_SCHEDULE.Enum(), UpdateMaintenanceSchedule: req}, ctx)
}

func (a *Client) StartMaintenance(req *Call_StartMaintenance, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_START_MAINTENANCE.Enum(), StartMaintenance: req}, ctx)
}

func (a *Client) StopMaintenance(req *Call_StopMaintenance, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_STOP_MAINTENANCE.Enum(), StopMaintenance: req}, ctx)
}

func (a *Client) GetQuota(ctx context.Context) (*Response_GetQuota, error) {
	resp, err := a.call(&Call{Type: Call_GET_QUOTA.Enum()}, ctx)
	return resp.GetQuota, err
}

func (a *Client) SetQuota(req *Call_SetQuota, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_SET_QUOTA.Enum(), SetQuota: req}, ctx)
}

func (a *Client) UpdateQuota(req *Call_RemoveQuota, ctx context.Context) error {
	return a.callNoResponse(&Call{Type: Call_REMOVE_QUOTA.Enum(), RemoveQuota: req}, ctx)
}

// returns pointer to empty struct in case of error
func (a *Client) call(c *Call, ctx context.Context) (*Response, error) {
	ev := &Response{}
	resp, err := a.c.Do(c, ctx)
	if err != nil {
		return ev, err
	}
	err = resp.Read(ev)
	resp.Close()
	return ev, err
}

func (a *Client) callNoResponse(c *Call, ctx context.Context) error {
	resp, err := a.c.Do(c, ctx)
	if resp != nil {
		resp.Close()
	}
	return err
}

func NewEventStream(c client.Client, ctx context.Context) <-chan *EventOrErr {
	events := make(chan *EventOrErr)
	go func() {
		resp, err := c.Do(&Call{Type: Call_SUBSCRIBE.Enum()}, ctx)
		for err == nil {
			msg := &Event{}
			err = resp.Read(msg)

			events <- &EventOrErr{Event: msg, Err: err}
		}

		if resp != nil {
			resp.Close()
		}
		events <- &EventOrErr{Err: err}
		close(events)
	}()

	return events
}
