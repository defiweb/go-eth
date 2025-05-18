package transport

/*
type Selector interface {
	Select(ctx context.Context, count int, method string, args ...any) (Transport, error)
}

type RoundRobinSelector struct {
	selected   int
	transports []Transport
}

func NewRoundRobinSelector(transports ...Transport) *RoundRobinSelector {
	return &RoundRobinSelector{transports: transports}
}

func (s *RoundRobinSelector) Select(_ context.Context, _ string, _ ...any) (Transport, error) {
	s.selected = (s.selected + 1) % len(s.transports)
	return s.transports[s.selected], nil
}

type Router struct {
	mu       sync.Mutex
	selector Selector
}

func NewRouter(selector Selector) *Router {
	return &Router{selector: selector}
}

// Call implements the Transport interface.
func (c *Router) Call(ctx context.Context, result any, method string, args ...any) error {
	transport, err := c.selectTransport(ctx, method, args...)
	if err != nil {
		return err
	}
	return transport.Call(ctx, result, method, args...)
}

func (c *Router) selectTransport(ctx context.Context, method string, args ...any) (Transport, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.selector.Select(ctx, method, args...)
}
*/
