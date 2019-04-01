package v1

import "strings"

func AddNotifier(es *[]Notifier, e *Notifier) bool {
	if !notifierInSlice(es, e) {
		*es = append(*es, *e)
		return true
	}
	return false
}

func RemoveNotifier(es *[]Notifier, e *Notifier) []Notifier {
	var newSubs []Notifier
	for _, b := range *es {
		if !b.Equal(e) {
			newSubs = append(newSubs, *e)
		}
	}
	return newSubs
}

func (s *Notifier) Equal(e *Notifier) bool {
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

func notifierInSlice(es *[]Notifier, e *Notifier) bool {
	for _, b := range *es {
		if b.Equal(e) {
			return true
		}
	}
	return false
}

func (n *Notifier) GetEventNotifier() EventNotifier {
	if n.Spec.Slack != nil {
		return n.Spec.Slack
	}
	return nil
}

func escapeString(str string) string {
	// Escape double quotes (")
	return strings.Replace(str, `"`, `\"`, -1)
}
