package v1

import (
	"encoding/json"
	"regexp"

	"github.com/nsf/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

func (s *EventSubscription) Subscribed(e *corev1.Event) (bool, error) {
	// Check if the event message is a match
	if s.Spec.MatchMessage != "" {
		matchMessage, err := regexp.MatchString(s.Spec.MatchMessage, e.Message)
		if err != nil {
			return false, err
		}
		if !matchMessage {
			return false, nil
		}
	}

	// Check if MatchReason is a match
	if s.Spec.MatchReason != "" {
		matchReason, err := regexp.MatchString(s.Spec.MatchReason, e.Reason)
		if err != nil {
			return false, err
		}
		if !matchReason {
			return false, nil
		}
	}

	// Check if MatchType is a match
	if s.Spec.MatchType != "" {
		matchType, err := regexp.MatchString(s.Spec.MatchType, e.Type)
		if err != nil {
			return false, err
		}
		if !matchType {
			return false, nil
		}
	}

	// Check if MatchObject is set and if so, compare it with InvolvedObject
	if (s.Spec.MatchObject != corev1.ObjectReference{}) {
		eventOut, err := json.Marshal(e.InvolvedObject.DeepCopy())
		if err != nil {
			return false, err
		}
		matchOut, err := json.Marshal(s.Spec.MatchObject.DeepCopy())
		if err != nil {
			return false, err
		}

		options := jsondiff.DefaultConsoleOptions()
		diff, _ := jsondiff.Compare(eventOut, matchOut, &options)
		if diff != jsondiff.FullMatch && diff != jsondiff.SupersetMatch {
			return false, nil
		}
	}

	// Nothing else to check, subscription must match
	return true, nil
}

func (s *EventSubscription) Equal(e *EventSubscription) bool {
	if s.TypeMeta.String() != e.TypeMeta.String() {
		return false
	}
	if s.ObjectMeta.GetName() != e.ObjectMeta.GetName() {
		return false
	}
	if s.ObjectMeta.GetNamespace() != e.ObjectMeta.GetNamespace() {
		return false
	}
	if s.Spec != e.Spec {
		return false
	}

	return true
}
