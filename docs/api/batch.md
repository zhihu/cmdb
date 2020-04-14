# Execute multiple requests atomically

**URL** : `/v1/api/batch`

A batch request consists of multiple API calls combined into on HTTP request, which can be sent to the endpoint and executed atomically.

## Format of a batch request

A batch request is a single standard HTTP request containing multiple API calls, using the `multipart/mixed` content-type. Within that
main HTTP request, each of the parts consists a nested HTTP request.

Each part begins with its own set of HTTP headers. After the server unwraps the batch request into separate requests, the part headers are ignored.

The body of each part is itself a complete HTTP request, with its own verb, URL, headers, and body. The HTTP request must only contain the path portion of the URL; full URLs are not allowed in batch requests.

The HTTP headers for the outer batch request, except for the Content- headers such as Content-Type, apply to every request in the batch. If you specify a given HTTP header in both the outer request and an individual call, then the individual call header's value overrides the outer batch request header's value. The headers for an individual call apply only to that call.

For example, if you provide an Authorization header for a specific call, then that header applies only to that call. If you provide an Authorization header for the outer request, then that header applies to all of the individual calls unless they override it with Authorization headers of their own.

When the server receives the batched request, it applies the outer request's query parameters and headers (as appropriate) to each part, and then treats each part as if it were a separate HTTP request.

## Response to a batch request

The server's response is a single standard HTTP response with a multipart/mixed content type; each part is the response to one of the requests in the batched request, in the same order as the requests.

Like the parts in the request, each response part contains a complete HTTP response, including a status code, headers, and body. And like the parts in the request, each response part is preceded by a Content-Type header that marks the beginning of the part.

## Notes
* The server perform your calls strictly in the order as the requests and they will execute atomically within the same transaction.
