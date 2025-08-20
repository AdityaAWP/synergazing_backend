package helper

// Timeline status constants
const (
	TimelineStatusNotStarted = "not-started"
	TimelineStatusInProgress = "in-progress"
	TimelineStatusDone       = "done"
)

// TimelineStatusOption represents a timeline status option for frontend selection
type TimelineStatusOption struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Color       string `json:"color,omitempty"`
}

// GetTimelineStatusOptions returns all available timeline status options
// This is used by the frontend to populate dropdown/select menus
func GetTimelineStatusOptions() []TimelineStatusOption {
	return []TimelineStatusOption{
		{
			Value:       TimelineStatusNotStarted,
			Label:       "Not Started",
			Description: "This timeline item hasn't been started yet",
			Color:       "#6B7280", // Gray
		},
		{
			Value:       TimelineStatusInProgress,
			Label:       "In Progress",
			Description: "This timeline item is currently being worked on",
			Color:       "#F59E0B", // Yellow/Orange
		},
		{
			Value:       TimelineStatusDone,
			Label:       "Done",
			Description: "This timeline item has been completed",
			Color:       "#10B981", // Green
		},
	}
}

// GetValidTimelineStatuses returns a slice of valid timeline status values
func GetValidTimelineStatuses() []string {
	options := GetTimelineStatusOptions()
	statuses := make([]string, len(options))
	for i, option := range options {
		statuses[i] = option.Value
	}
	return statuses
}

// IsValidTimelineStatus checks if a given status is valid
func IsValidTimelineStatus(status string) bool {
	validStatuses := GetValidTimelineStatuses()
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// GetDefaultTimelineStatus returns the default timeline status
func GetDefaultTimelineStatus() string {
	return TimelineStatusNotStarted
}

// GetTimelineStatusByValue returns the timeline status option for a given value
func GetTimelineStatusByValue(value string) *TimelineStatusOption {
	options := GetTimelineStatusOptions()
	for _, option := range options {
		if option.Value == value {
			return &option
		}
	}
	return nil
}
