package minrevpro_test

import (
	"reflect"
	"testing"
)

type mustable struct {
	t *testing.T
}

func (m mustable) must() {
	m.t.Helper()
	if m.t.Failed() {
		m.t.FailNow()
	}
}

// equal string

func assertEqualString(t *testing.T, want, got string, msg string, args ...interface{}) mustable {
	t.Helper()
	check := want == got
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// equal int

func assertEqualInt(t *testing.T, want, got int, msg string, args ...interface{}) mustable {
	t.Helper()
	check := want == got
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// true

func assertTrue(t *testing.T, got bool, msg string, args ...interface{}) mustable {
	t.Helper()
	check := got == true
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// false

func assertFalse(t *testing.T, got bool, msg string, args ...interface{}) mustable {
	t.Helper()
	check := got == false
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// notNil

func assertNotNil(t *testing.T, got interface{}, msg string, args ...interface{}) mustable {
	t.Helper()
	check := !(got == nil || reflect.ValueOf(got).IsNil())
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// nil

func assertNil(t *testing.T, got interface{}, msg string, args ...interface{}) mustable {
	t.Helper()
	check := got == nil || reflect.ValueOf(got).IsNil()
	failIfNOK(t, check, msg, args...)
	return mustable{t}
}

// helpers

func failIfNOK(t *testing.T, check bool, msg string, args ...interface{}) {
	if check {
		return
	}
	t.Helper()
	t.Errorf(msg, args...)
}
