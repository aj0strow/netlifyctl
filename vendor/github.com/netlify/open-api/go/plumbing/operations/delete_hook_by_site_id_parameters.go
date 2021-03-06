package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteHookBySiteIDParams creates a new DeleteHookBySiteIDParams object
// with the default values initialized.
func NewDeleteHookBySiteIDParams() *DeleteHookBySiteIDParams {
	var ()
	return &DeleteHookBySiteIDParams{}
}

/*DeleteHookBySiteIDParams contains all the parameters to send to the API endpoint
for the delete hook by site Id operation typically these are written to a http.Request
*/
type DeleteHookBySiteIDParams struct {

	/*HookID*/
	HookID string
}

// WithHookID adds the hookId to the delete hook by site Id params
func (o *DeleteHookBySiteIDParams) WithHookID(HookID string) *DeleteHookBySiteIDParams {
	o.HookID = HookID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteHookBySiteIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	// path param hook_id
	if err := r.SetPathParam("hook_id", o.HookID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
