package skelet

// Authentication confirms that users are who they say they are.
// Authorization gives those users permission to access a resource.

// Middlewares order.
//	All: LimitMiddleware: Check if the request is not exceeding the limit. | Can abort before calling next.
// 		UserGroup: JwtMiddleware: Check JWT. Set user to context. | Can abort before calling next.
//			All: EndpointMiddleware: Executes endpoint. | Should not abort.
//		UserGroup: JwtMiddleware: NOOP
//	All: LimitMiddleware: Check aborted | Create access log. Maybe asynchronously, for not slowing down the requests. Write in batches.
