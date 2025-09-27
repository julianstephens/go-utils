package response

func New() *Responder {
	return &Responder{
		Encoder: NewJSONEncoder(),
		Before:  DefaultBefore,
		After:   DefaultAfter,
		OnError: DefaultOnError,
	}
}

// NewEmpty creates a new Responder with JSON encoder and no hooks.
func NewEmpty() *Responder {
	return &Responder{
		Encoder: NewJSONEncoder(),
		Before:  nil,
		After:   nil,
		OnError: nil,
	}
}

// NewWithLogging creates a new Responder with JSON encoder and logging hooks.
// This configuration logs requests, responses, and errors.
func NewWithLogging() *Responder {
	return &Responder{
		Encoder: NewJSONEncoder(),
		Before:  LoggingBefore,
		After:   LoggingAfter,
		OnError: LoggingOnError,
	}
}

// NewCustom creates a new Responder with custom configuration.
func NewCustom(encoder Encoder, before BeforeFunc, after AfterFunc, onError OnErrorFunc) *Responder {
	if encoder == nil {
		encoder = NewJSONEncoder()
	}

	return &Responder{
		Encoder: encoder,
		Before:  before,
		After:   after,
		OnError: onError,
	}
}
