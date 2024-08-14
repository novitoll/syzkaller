package main

import (
	"fmt"
	"testing"

	"github.com/google/syzkaller/pkg/report"
	"github.com/stretchr/testify/assert"
)

func Test_Crash_FullTitle(t *testing.T) {
	crash := &Crash{}

	tests := []struct {
		name          string
		report        *report.Report
		fromDashboard bool
		fromHub       bool
		expected      string
	}{
		{
			name:     "report title is filled",
			report:   &report.Report{Title: "foo"},
			expected: "foo",
		},
		{
			name:          "report title fromDashboard",
			report:        &report.Report{},
			fromDashboard: true,
			expected:      fmt.Sprintf("dashboard crash %p", crash),
		},
		{
			name:     "report title fromHub",
			report:   &report.Report{},
			fromHub:  true,
			expected: fmt.Sprintf("crash from hub %p", crash),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.report.Title == "" && !tt.fromDashboard && !tt.fromHub {
				assert.Panics(t, func() { crash.FullTitle() })
			} else {
				crash.Report = tt.report
				crash.fromDashboard = tt.fromDashboard
				crash.fromHub = tt.fromHub

				title := crash.FullTitle()
				assert.Equal(t, tt.expected, title)
			}
		})
	}
}
