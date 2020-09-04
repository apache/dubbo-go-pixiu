package context

type FilterFunc func(Context)

type FilterChain []FilterFunc
