// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package checkers_test

import (
	"testing"
	"time"

	gc "gopkg.in/check.v1"

	jc "github.com/juju/testing/checkers"
)

func Test(t *testing.T) { gc.TestingT(t) }

type CheckerSuite struct{}

var _ = gc.Suite(&CheckerSuite{})

func (s *CheckerSuite) TestHasPrefix(c *gc.C) {
	c.Assert("foo bar", jc.HasPrefix, "foo")
	c.Assert("foo bar", gc.Not(jc.HasPrefix), "omg")
}

func (s *CheckerSuite) TestHasSuffix(c *gc.C) {
	c.Assert("foo bar", jc.HasSuffix, "bar")
	c.Assert("foo bar", gc.Not(jc.HasSuffix), "omg")
}

func (s *CheckerSuite) TestContains(c *gc.C) {
	c.Assert("foo bar baz", jc.Contains, "foo")
	c.Assert("foo bar baz", jc.Contains, "bar")
	c.Assert("foo bar baz", jc.Contains, "baz")
	c.Assert("foo bar baz", gc.Not(jc.Contains), "omg")
}

func (s *CheckerSuite) TestTimeBetween(c *gc.C) {
	now := time.Now()
	earlier := now.Add(-1 * time.Second)
	later := now.Add(time.Second)

	check := func(value interface{}, start, end time.Time) (bool, string) {
		checker := jc.TimeBetween(start, end)
		return checker.Check([]interface{}{value}, nil)
	}

	value, msg := check(now, earlier, later)
	c.Assert(value, jc.IsTrue)
	c.Assert(msg, gc.Equals, "")
	// Later can be before earlier...
	value, msg = check(now, later, earlier)
	c.Assert(value, jc.IsTrue)
	c.Assert(msg, gc.Equals, "")

	value, msg = check(earlier, now, later)
	c.Assert(value, jc.IsFalse)
	c.Assert(msg, gc.Matches, `obtained time .* is before start time .*`)

	value, msg = check(later, now, earlier)
	c.Assert(value, jc.IsFalse)
	c.Assert(msg, gc.Matches, `obtained time .* is after end time .*`)

	value, msg = check(42, now, earlier)
	c.Assert(value, jc.IsFalse)
	c.Assert(msg, gc.Matches, `obtained value type must be time.Time`)

	// equality checking
	value, msg = check(earlier, earlier, later)
	c.Assert(value, jc.IsTrue)
	c.Assert(msg, gc.Equals, "")
	value, msg = check(later, earlier, later)
	c.Assert(value, jc.IsTrue)
	c.Assert(msg, gc.Equals, "")
}

func (s *CheckerSuite) TestSameContents(c *gc.C) {
	//// positive cases ////

	// same
	c.Check(
		[]int{1, 2, 3}, jc.SameContents,
		[]int{1, 2, 3})

	// empty
	c.Check(
		[]int{}, jc.SameContents,
		[]int{})

	// single
	c.Check(
		[]int{1}, jc.SameContents,
		[]int{1})

	// different order
	c.Check(
		[]int{1, 2, 3}, jc.SameContents,
		[]int{3, 2, 1})

	// multiple copies of same
	c.Check(
		[]int{1, 1, 2}, jc.SameContents,
		[]int{2, 1, 1})

	type test struct {
		s string
		i int
	}

	// test structs
	c.Check(
		[]test{{"a", 1}, {"b", 2}}, jc.SameContents,
		[]test{{"b", 2}, {"a", 1}})

	//// negative cases ////

	// different contents
	c.Check(
		[]int{1, 3, 2, 5}, gc.Not(jc.SameContents),
		[]int{5, 2, 3, 4})

	// different size slices
	c.Check(
		[]int{1, 2, 3}, gc.Not(jc.SameContents),
		[]int{1, 2})

	// different counts of same items
	c.Check(
		[]int{1, 1, 2}, gc.Not(jc.SameContents),
		[]int{1, 2, 2})

	/// Error cases ///
	//  note: for these tests, we can't use gc.Not, since Not passes the error value through
	// and checks with a non-empty error always count as failed
	// Oddly, there doesn't seem to actually be a way to check for an error from a Checker.

	// different type
	res, err := jc.SameContents.Check([]interface{}{
		[]string{"1", "2"},
		[]int{1, 2},
	}, []string{})
	c.Check(res, jc.IsFalse)
	c.Check(err, gc.Not(gc.Equals), "")

	// obtained not a slice
	res, err = jc.SameContents.Check([]interface{}{
		"test",
		[]int{1},
	}, []string{})
	c.Check(res, jc.IsFalse)
	c.Check(err, gc.Not(gc.Equals), "")

	// expected not a slice
	res, err = jc.SameContents.Check([]interface{}{
		[]int{1},
		"test",
	}, []string{})
	c.Check(res, jc.IsFalse)
	c.Check(err, gc.Not(gc.Equals), "")
}
