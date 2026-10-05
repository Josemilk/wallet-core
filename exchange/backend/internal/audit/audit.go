package audit

import (
    "context"
    "errors"
)

type Event struct { ID, ActorID, Action, Resource, RequestID, MetadataJSON string }
type Sink interface { Append(context.Context, Event) error }
func Validate(e Event) error { if e.ID==""||e.ActorID==""||e.Action==""||e.Resource==""||e.RequestID=="" { return errors.New("incomplete audit event") }; return nil }
