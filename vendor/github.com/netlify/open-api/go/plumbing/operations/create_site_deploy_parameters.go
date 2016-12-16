package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/netlify/open-api/go/models"
)

// NewCreateSiteDeployParams creates a new CreateSiteDeployParams object
// with the default values initialized.
func NewCreateSiteDeployParams() *CreateSiteDeployParams {
	var ()
	return &CreateSiteDeployParams{}
}

/*CreateSiteDeployParams contains all the parameters to send to the API endpoint
for the create site deploy operation typically these are written to a http.Request
*/
type CreateSiteDeployParams struct {

	/*Deploy*/
	Deploy *models.DeployFiles
	/*SiteID*/
	SiteID string
}

// WithDeploy adds the deploy to the create site deploy params
func (o *CreateSiteDeployParams) WithDeploy(Deploy *models.DeployFiles) *CreateSiteDeployParams {
	o.Deploy = Deploy
	return o
}

// WithSiteID adds the siteId to the create site deploy params
func (o *CreateSiteDeployParams) WithSiteID(SiteID string) *CreateSiteDeployParams {
	o.SiteID = SiteID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *CreateSiteDeployParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	if o.Deploy == nil {
		o.Deploy = new(models.DeployFiles)
	}

	if err := r.SetBodyParam(o.Deploy); err != nil {
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
