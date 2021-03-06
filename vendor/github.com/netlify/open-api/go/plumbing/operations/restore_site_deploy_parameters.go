package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// NewRestoreSiteDeployParams creates a new RestoreSiteDeployParams object
// with the default values initialized.
func NewRestoreSiteDeployParams() *RestoreSiteDeployParams {
	var ()
	return &RestoreSiteDeployParams{}
}

/*RestoreSiteDeployParams contains all the parameters to send to the API endpoint
for the restore site deploy operation typically these are written to a http.Request
*/
type RestoreSiteDeployParams struct {

	/*DeployID*/
	DeployID string
	/*SiteID*/
	SiteID string
}

// WithDeployID adds the deployId to the restore site deploy params
func (o *RestoreSiteDeployParams) WithDeployID(DeployID string) *RestoreSiteDeployParams {
	o.DeployID = DeployID
	return o
}

// WithSiteID adds the siteId to the restore site deploy params
func (o *RestoreSiteDeployParams) WithSiteID(SiteID string) *RestoreSiteDeployParams {
	o.SiteID = SiteID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *RestoreSiteDeployParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	// path param deploy_id
	if err := r.SetPathParam("deploy_id", o.DeployID); err != nil {
		return err
	}

	// path param site_id
	if err := r.SetPathParam("site_id", o.SiteID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
