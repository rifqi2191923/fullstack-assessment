package repository

import (
	"testing"

	"github.com/example/fullstack-assessment/backend/internal/model"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		status  string
		wantErr bool
	}{
		{
			name:    "valid todo",
			title:   "Test task",
			status:  "todo",
			wantErr: false,
		},
		{
			name:    "valid in progress",
			title:   "Test task",
			status:  "in_progress",
			wantErr: false,
		},
		{
			name:    "valid done",
			title:   "Test task",
			status:  "done",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			status:  "todo",
			wantErr: true,
		},
		{
			name:    "invalid status",
			title:   "Test task",
			status:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInput(
				model.TaskInput{
					Title:  tt.title,
					Status: tt.status,
				},
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"ValidateInput() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}