package utils

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	case diff < 12*30*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(diff.Hours()/24/30))
	default:
		return fmt.Sprintf("%d years ago", int(diff.Hours()/24/365))
	}
}
