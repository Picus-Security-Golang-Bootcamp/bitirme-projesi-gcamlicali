// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ItemQuantity item quantity
//
// swagger:model ItemQuantity
type ItemQuantity struct {

	// quantity
	// Required: true
	Quantity *int32 `json:"quantity"`
}

// Validate validates this item quantity
func (m *ItemQuantity) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateQuantity(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ItemQuantity) validateQuantity(formats strfmt.Registry) error {

	if err := validate.Required("quantity", "body", m.Quantity); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this item quantity based on context it is used
func (m *ItemQuantity) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ItemQuantity) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ItemQuantity) UnmarshalBinary(b []byte) error {
	var res ItemQuantity
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
