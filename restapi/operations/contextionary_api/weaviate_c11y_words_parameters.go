/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2019 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * DESIGN & CONCEPT: Bob van Luijt (@bobvanluijt)
 * CONTACT: hello@creativesoftwarefdn.org
 */ // Code generated by go-swagger; DO NOT EDIT.

package contextionary_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewWeaviateC11yWordsParams creates a new WeaviateC11yWordsParams object
// no default values defined in spec.
func NewWeaviateC11yWordsParams() WeaviateC11yWordsParams {

	return WeaviateC11yWordsParams{}
}

// WeaviateC11yWordsParams contains all the bound params for the weaviate c11y words operation
// typically these are obtained from a http.Request
//
// swagger:parameters weaviate.c11y.words
type WeaviateC11yWordsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*CamelCase list of words to validate.
	  Required: true
	  In: path
	*/
	Words string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewWeaviateC11yWordsParams() beforehand.
func (o *WeaviateC11yWordsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rWords, rhkWords, _ := route.Params.GetOK("words")
	if err := o.bindWords(rWords, rhkWords, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindWords binds and validates parameter Words from path.
func (o *WeaviateC11yWordsParams) bindWords(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Words = raw

	return nil
}