package context

// FilterFunc filter func, filter
type FilterFunc func(Context)

// FilterChain filter chain
type FilterChain []FilterFunc
