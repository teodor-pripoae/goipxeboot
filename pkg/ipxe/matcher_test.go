package ipxe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticIPMatcherInvalid(t *testing.T) {
	matcher, err := MatchStaticIP("invalid")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "invalid IP", err.Error(), "error message should be correct")

	matcher, err = MatchStaticIP("300.300.300.300")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "invalid IP", err.Error(), "error message should be correct")

	matcher, err = MatchStaticIP("2001:0000:130F:0000:0000:09C0:876A:130B")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "IPv6 not supported", err.Error(), "error message should be correct")

	matcher, err = MatchStaticIP("192.168.122.1/24")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "invalid IP", err.Error(), "error message should be correct")
}

func TestStaticIPMatcher(t *testing.T) {
	matcher, err := MatchStaticIP("192.168.122.100")
	assert.NoError(t, err, "no error should be returned")

	assert.True(t, matcher("192.168.122.100"), "should match")
	assert.False(t, matcher("192.168.1.100"), "should not match")
}

func TestIPRangeInvalid(t *testing.T) {
	matcher, err := MatchIPRange("invalid")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "invalid CIDR address: invalid", err.Error(), "error message should be correct")

	matcher, err = MatchIPRange("10.20.30.40-50")
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "invalid CIDR address: 10.20.30.40-50", err.Error(), "error message should be correct")
}

func TestIPRangeOK(t *testing.T) {
	matcher, err := MatchIPRange("192.168.122.1/24")
	assert.NoError(t, err, "no error should be returned")
	assert.NotNil(t, matcher, "matcher should not be nil")

	assert.True(t, matcher("192.168.122.100"), "should match")
	assert.True(t, matcher("192.168.122.1"), "should match")
	assert.True(t, matcher("192.168.122.243"), "should match")

	assert.False(t, matcher("192.168.1.100"), "should not match")
}

func TestDetectNoMatcher(t *testing.T) {
	matcher, err := DetectMatcher("invalid", DefaultMatchers)
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "no matcher found for invalid", err.Error(), "error message should be correct")

	matcher, err = DetectMatcher("10.20.30.40-50", DefaultMatchers)
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "no matcher found for 10.20.30.40-50", err.Error(), "error message should be correct")

	matcher, err = DetectMatcher("2001:0000:130F:0000:0000:09C0:876A:130B", DefaultMatchers)
	assert.Error(t, err, "error should be returned")
	assert.Nil(t, matcher, "matcher should be nil")
	assert.Equal(t, "no matcher found for 2001:0000:130F:0000:0000:09C0:876A:130B", err.Error(), "error message should be correct")
}

func TestDetectMatcher(t *testing.T) {
	matcher, err := DetectMatcher("192.168.122.100", DefaultMatchers)
	assert.NoError(t, err, "no error should be returned")
	assert.NotNil(t, matcher, "matcher should not be nil")
	assert.True(t, matcher("192.168.122.100"), "should match")
	assert.False(t, matcher("192.168.122.1"), "should not match")

	matcher, err = DetectMatcher("192.168.122.1/24", DefaultMatchers)
	assert.NoError(t, err, "no error should be returned")
	assert.NotNil(t, matcher, "matcher should not be nil")
	assert.True(t, matcher("192.168.122.100"), "should match")
	assert.True(t, matcher("192.168.122.1"), "should match")
	assert.False(t, matcher("192.168.1.100"), "should not match")
}
