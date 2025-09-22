// Package request contains request-parsing helpers for HTTP handlers used
// by the project and its tests.
//
// It provides small utilities to parse query parameters, form values, and
// decode JSON request bodies with sensible defaults and error handling.
//
// # Common helpers
//
//   - QueryInt(r *http.Request, key string, defaultVal int) int
//     Parse an integer query value returning a default when missing or invalid.
//
//   - QueryBool(r *http.Request, key string, defaultVal bool) bool
//     Parse a boolean query parameter with a default.
//
//   - DecodeJSONBody(ctx context.Context, r *http.Request, v interface{}) error
//     Decode a JSON body into v and return helpful error messages for clients.
//
// Example usage in a handler
//
//	func handleListUsers(w http.ResponseWriter, r *http.Request) {
//	    page := request.QueryInt(r, "page", 1)
//	    perPage := request.QueryInt(r, "per_page", 25)
//	    // use page/perPage to fetch results
//	}
//
// The helpers are intentionally small and composable, intended to reduce
// repetition in handlers and make tests easier to write.
package request
