package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/google/syzkaller/pkg/corpus"
	"github.com/google/syzkaller/pkg/mgrconfig"
	"github.com/google/syzkaller/pkg/report"
	"github.com/google/syzkaller/sys/targets"
	"github.com/google/syzkaller/vm"
	"github.com/google/syzkaller/vm/vmimpl"
	"github.com/google/syzkaller/vm/vmtest"
)

func initVm() {
	ctor := func(env *vmimpl.Env) (vmimpl.Pool, error) {
		return &vmtest.TestPool{}, nil
	}
	vmimpl.Register("test", vmimpl.Type{
		Ctor:        ctor,
		Preemptible: true,
	})
}

func TestRunManager(t *testing.T) {
	initVm()

	tests := []struct {
		name           string
		mode           Mode
		cfg            *mgrconfig.Config
		isFlagBenchSet bool
	}{
		{
			name: "ok",
			cfg:  &mgrconfig.Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.Type = "test"
			tt.cfg.SysTarget = targets.Get(targets.Linux, targets.AMD64)
			vmPool, err := vm.Create(tt.cfg, *flagDebug)
			assert.NoError(t, err)

			mgr := newManager(tt.cfg, tt.mode, vmPool, nil, nil)
			impl := new(managerImplMock)
			shouldPreloadCorpus := mgr.mode == ModeFuzzing || mgr.mode == ModeCorpusTriage || mgr.mode == ModeCorpusRun

			impl.On("initStats")
			if shouldPreloadCorpus {
				impl.On("preloadCorpus")
			}
			impl.On("initHTTP").
				On("corpusInputHandler", mock.Anything).
				On("trackUsedFiles").
				On("initRPC")

			if tt.cfg.DashboardAddr != "" {
				impl.On("initDash")
			}

			if tt.cfg.AssetStorage.IsEmpty() {
				impl.On("initAssetStorage")
			}

			if tt.isFlagBenchSet {
				impl.On("initBench")
			}

			impl.On("heartbeatLoop")

			mgr.run(impl)

			if !shouldPreloadCorpus {
				// mgr.corpusPreload channels should be closed
				_, ok := <-mgr.corpusPreload
				assert.False(t, ok, "corpusPreload should be closed")
			}
		})
	}
}

func TestCrash_FullTitle(t *testing.T) {
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

type managerImplMock struct {
	mock.Mock
}

func (m *managerImplMock) initStats() {
	_ = m.Called()
}

func (m *managerImplMock) initHTTP() {
	_ = m.Called()
}

func (m *managerImplMock) initRPC() {
	_ = m.Called()
}

func (m *managerImplMock) initDash() {
	_ = m.Called()
}

func (m *managerImplMock) initAssetStorage() {
	_ = m.Called()
}

func (m *managerImplMock) initBench() {
	_ = m.Called()
}

func (m *managerImplMock) heartbeatLoop() {
	_ = m.Called()
}

func (m *managerImplMock) preloadCorpus() {
	_ = m.Called()
}

func (m *managerImplMock) corpusInputHandler(corpusUpdates <-chan corpus.NewItemEvent) {
	_ = m.Called(corpusUpdates)
}

func (m *managerImplMock) trackUsedFiles() {
	_ = m.Called()
}

func (m *managerImplMock) processFuzzingResults(ctx context.Context) {
	_ = m.Called(ctx)
}
