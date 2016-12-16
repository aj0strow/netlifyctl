package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/netlify/open-api/go/models"
)

// NewUpdateSiteParams creates a new UpdateSiteParams object
// with the default values initialized.
func NewUpdateSiteParams() *UpdateSiteParams {
	var ()
	return &UpdateSiteParams{}
}

/*UpdateSiteParams contains all the parameters to send to the API endpoint
for the update site operation typically these are written to a http.Request
*/
type UpdateSiteParams struct {

	/*Site*/
	Site *models.Site
	/*SiteID*/
	SiteID string
}

// WithSite adds the site to the update site params
func (o *UpdateSiteParams) WithSite(Site *models.Site) *UpdateSiteParams {
	o.Site = Site
	return o
}

// WithSiteID adds the siteId to the update site params
func (o *UpdateSiteParams) WithSiteID(SiteID string) *UpdateSiteParams {
	o.SiteID = SiteID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateSiteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	if o.Site == nil {
		o.Site = new(models.Site)
	}

	if err := r.SetBodyParam(o.Site); err != nil {
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
