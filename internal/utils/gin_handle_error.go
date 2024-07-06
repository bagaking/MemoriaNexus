package utils

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/khicago/irr"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse defines the standard error response structure.
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// GinHandleErrorOption defines the type for functional options.
type GinHandleErrorOption func(*ginHandleErrorOptions)

// ginHandleErrorOptions holds the options for GinHandleError.
type ginHandleErrorOptions struct {
	requestContents map[string]any
	extra           map[string]any
}

// GinErrWithReqContents sets the requestContents option.
func GinErrWithReqContents(params map[string]any) GinHandleErrorOption {
	return func(opts *ginHandleErrorOptions) {
		opts.requestContents = params
	}
}

// GinErrWithExtra sets the requestContents option.
func GinErrWithExtra(key string, val any) GinHandleErrorOption {
	return func(opts *ginHandleErrorOptions) {
		if opts.extra == nil {
			opts.extra = make(map[string]any)
		}
		opts.extra[key] = val
	}
}

// GinErrWithReqBody sets the requestBody
func GinErrWithReqBody(req any) GinHandleErrorOption {
	return func(opts *ginHandleErrorOptions) {
		if opts.requestContents == nil {
			opts.requestContents = map[string]any{}
		}
		opts.requestContents["_req_body"] = req
	}
}

// GinHandleError handles errors by logging them and sending a ErrorResponse.
func GinHandleError(c *gin.Context, log logrus.FieldLogger, status int, err error, msg string, options ...GinHandleErrorOption) {
	resp := ErrorResponse{
		Message: msg,
	}
	if err != nil {
		resp.Error = err.Error() // Ensure the error is converted to a string
	}

	// Initialize options with defaults
	opts := &ginHandleErrorOptions{
		requestContents: make(map[string]any),
	}

	// Apply options
	for _, opt := range options {
		opt(opts)
	}

	// Collect basic parameters from the request context
	queries := make(map[string]any)
	params := make(map[string]any)

	// Collect query parameters
	for key, values := range c.Request.URL.Query() {
		if len(values) < 1 {
			continue
		}
		if len(values) == 1 {
			queries[key] = values[0]
		} else if len(values) > 1 {
			queries[key] = values
		}
	}

	// Collect URL parameters
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}

	// Add the error, message, and stack trace to the log fields
	fields := logrus.Fields{
		"error":      err,
		"message":    msg,
		"queries":    queries,
		"params":     params,
		"contents":   opts.requestContents,
		"extra":      opts.extra,
		"stacktrace": getGinSimplifiedStackTrace(),
	}

	var ir irr.IRR
	if errors.As(err, &ir) {
		fields["irr.track"] = ir.ToString(true, " ++ ")
	}

	// Log the error with different levels based on the status code
	switch {
	case status >= http.StatusInternalServerError:
		log.WithFields(fields).Errorf("internal server error") // Internal server errors are logged as errors
	case status >= http.StatusBadRequest:
		log.WithFields(fields).Warnf("Client error; stacktrace") // Client errors are logged as warnings
	default:
		log.WithFields(fields).Infof("Other case; stacktrace") // Other cases (e.g., redirects) are logged as info
	}

	c.JSON(status, resp)
}

// getGinSimplifiedStackTrace returns a simplified stack trace.
func getGinSimplifiedStackTrace() string {
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stackTrace []string
	for i := 0; i < 100; i++ {
		frame, more := frames.Next()
		// Extract the package, function, file, and line number
		parts := strings.Split(frame.Function, "/")
		funcName := parts[len(parts)-1]
		fileLine := fmt.Sprintf("%s:%d", frame.File, frame.Line)

		// Simplify known patterns
		switch {
		case strings.Contains(funcName, "gin.(*Engine).ServeHTTP") || strings.Contains(funcName, "gin.(*Engine).handleHTTPRequest"):
			stackTrace = append(stackTrace, "[[ServeHTTP by gin]]")
			more = false
		case strings.Contains(funcName, "gin.(*Context).Next"):
			// Skip repetitive gin context next calls
			continue
		case strings.Contains(funcName, "gin"):
			stackTrace = append(stackTrace, fmt.Sprintf(" -NEXT- %s", funcName))
		case strings.Contains(funcName, "utils.GinHandleError"):
			stackTrace = append(stackTrace, " -CALL HandleError-")
		case strings.Contains(funcName, "runtime."):
			stackTrace = append(stackTrace, funcName)
		default:
			stackTrace = append(stackTrace, fmt.Sprintf(" -FUNC- ┃%s┃ (%s)", funcName, fileLine))
		}

		if !more {
			break
		}
	}

	// Reverse the stack trace for readability
	for i, j := 0, len(stackTrace)-1; i < j; i, j = i+1, j-1 {
		stackTrace[i], stackTrace[j] = stackTrace[j], stackTrace[i]
	}

	// Join the stack trace elements with new lines
	return strings.Join(stackTrace, "\n")
}
