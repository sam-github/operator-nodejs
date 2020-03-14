package kwrap

import (
	"regexp"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type NsPred []string

var _ predicate.Predicate = NsPred(nil)

func MatchNamespacePredicate(watchNamespace string) NsPred {
	sep := regexp.MustCompile(`/[\s,]*([^\s]+)[\s,]*/`)
	ns := sep.FindAllString(watchNamespace, 0)
	var watch []string
	for _, ns := range ns {
		if ns != "" {
			watch = append(watch, ns)
		}
	}
	if watch == nil {
		watch = append(watch, "*")
	}

	return NsPred(watch)
}

func (p NsPred) watching(ns string) bool {
	for _, w := range p {
		if w == "*" || w == ns {
			return true
		}
	}
	return false
}

func (p NsPred) Create(e event.CreateEvent) bool {
	return p.watching(e.Meta.GetNamespace())
}

func (p NsPred) Delete(e event.DeleteEvent) bool {
	return p.watching(e.Meta.GetNamespace())
}

func (p NsPred) Update(e event.UpdateEvent) bool {
	return p.watching(e.MetaOld.GetNamespace())
}

func (p NsPred) Generic(e event.GenericEvent) bool {
	return p.watching(e.Meta.GetNamespace())
}
