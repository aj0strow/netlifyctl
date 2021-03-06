package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// NewProvisionSiteTLSCertificateParams creates a new ProvisionSiteTLSCertificateParams object
// with the default values initialized.
func NewProvisionSiteTLSCertificateParams() *ProvisionSiteTLSCertificateParams {
	var ()
	return &ProvisionSiteTLSCertificateParams{}
}

/*ProvisionSiteTLSCertificateParams contains all the parameters to send to the API endpoint
for the provision site TLS certificate operation typically these are written to a http.Request
*/
type ProvisionSiteTLSCertificateParams struct {

	/*CaCertificates*/
	CaCertificates *string
	/*Certificate*/
	Certificate *string
	/*Key*/
	Key *string
	/*SiteID*/
	SiteID string
}

// WithCaCertificates adds the caCertificates to the provision site TLS certificate params
func (o *ProvisionSiteTLSCertificateParams) WithCaCertificates(CaCertificates *string) *ProvisionSiteTLSCertificateParams {
	o.CaCertificates = CaCertificates
	return o
}

// WithCertificate adds the certificate to the provision site TLS certificate params
func (o *ProvisionSiteTLSCertificateParams) WithCertificate(Certificate *string) *ProvisionSiteTLSCertificateParams {
	o.Certificate = Certificate
	return o
}

// WithKey adds the key to the provision site TLS certificate params
func (o *ProvisionSiteTLSCertificateParams) WithKey(Key *string) *ProvisionSiteTLSCertificateParams {
	o.Key = Key
	return o
}

// WithSiteID adds the siteId to the provision site TLS certificate params
func (o *ProvisionSiteTLSCertificateParams) WithSiteID(SiteID string) *ProvisionSiteTLSCertificateParams {
	o.SiteID = SiteID
	return o
}

// WriteToRequest writes these params to a swagger request
func (o *ProvisionSiteTLSCertificateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	var res []error

	if o.CaCertificates != nil {

		// query param ca_certificates
		var qrCaCertificates string
		if o.CaCertificates != nil {
			qrCaCertificates = *o.CaCertificates
		}
		qCaCertificates := qrCaCertificates
		if qCaCertificates != "" {
			if err := r.SetQueryParam("ca_certificates", qCaCertificates); err != nil {
				return err
			}
		}

	}

	if o.Certificate != nil {

		// query param certificate
		var qrCertificate string
		if o.Certificate != nil {
			qrCertificate = *o.Certificate
		}
		qCertificate := qrCertificate
		if qCertificate != "" {
			if err := r.SetQueryParam("certificate", qCertificate); err != nil {
				return err
			}
		}

	}

	if o.Key != nil {

		// query param key
		var qrKey string
		if o.Key != nil {
			qrKey = *o.Key
		}
		qKey := qrKey
		if qKey != "" {
			if err := r.SetQueryParam("key", qKey); err != nil {
				return err
			}
		}

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
