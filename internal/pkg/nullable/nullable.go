package nullable

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Uint struct {
	Value *uint
	Set   bool
}

func (n *Uint) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v uint
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.Value = &v
	return nil
}

type String struct {
	Value *string
	Set   bool
}

func (n *String) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.Value = &v
	return nil
}

type UUID struct {
	Value *uuid.UUID
	Set   bool
}

func (n *UUID) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	parsed, err := uuid.Parse(v)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	n.Value = &parsed
	return nil
}

type Time struct {
	Value *time.Time
	Set   bool
}

func (n *Time) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	n.Value = &parsed
	return nil
}

type Bool struct {
	Value *bool
	Set   bool
}

func (n *Bool) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	n.Value = &parsed
	return nil
}
