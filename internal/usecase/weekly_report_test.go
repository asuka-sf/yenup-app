package usecase

import (
	"context"
	"errors"
	"testing"

	"yenup/internal/domain/rate"

	"github.com/stretchr/testify/assert"
)

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name          string
		mockRates     []*rate.Rate
		mockReadErr   error
		mockNotifyErr error
		wantErr       bool
	}{
		{
			name:          "success",
			mockRates:     testValidRates,
			mockReadErr:   nil,
			mockNotifyErr: nil,
			wantErr:       false,
		},
		{
			name:          "err: cannot load a JSON file",
			mockRates:     nil,
			mockReadErr:   errors.New("failed to load a JSON file"),
			mockNotifyErr: nil,
			wantErr:       true,
		},
		{
			name:          "err: rates is empty",
			mockRates:     nil,
			mockReadErr:   nil,
			mockNotifyErr: nil,
			wantErr:       true,
		},
		{
			name:          "err: duplicated dates",
			mockRates:     testDuplicatedDate,
			mockReadErr:   nil,
			mockNotifyErr: nil,
			wantErr:       true,
		},
		{
			name:          "err: inconsistent base",
			mockRates:     testInconsistentBase,
			mockReadErr:   nil,
			mockNotifyErr: nil,
			wantErr:       true,
		},
		{
			name:          "err: inconsistent target",
			mockRates:     testInconsistentTarget,
			mockReadErr:   nil,
			mockNotifyErr: nil,
			wantErr:       true,
		},
		{
			name:          "err: fail to notify",
			mockRates:     testValidRates,
			mockReadErr:   nil,
			mockNotifyErr: errors.New("failed to notify"),
			wantErr:       true,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockStorageClient{
				rates:   tt.mockRates,
				readErr: tt.mockReadErr,
			}
			notifier := &MockNotifier{err: tt.mockNotifyErr}
			uc := NewWeeklyReporter(storage, notifier)
			err := uc.GenerateReport(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !tt.wantErr {
				assert.Contains(t, notifier.msg, "Average")
				assert.Contains(t, notifier.msg, "Max")
				assert.Contains(t, notifier.msg, "Min")
			}
		})
	}

}
