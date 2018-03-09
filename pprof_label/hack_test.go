package pprof_label

import (
	"reflect"
	"testing"
)

func TestSetGoroutineLabels(t *testing.T) {
	sync := make(chan struct{})

	wantLabels := LabelMap{}

	if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
		t.Errorf("Expected parent goroutine's profile labels to be empty before test, got %v", gotLabels)
	}
	go func() {
		if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
			t.Errorf("Expected child goroutine's profile labels to be empty before test, got %v", gotLabels)
		}
		sync <- struct{}{}
	}()
	<-sync

	wantLabels = map[string]string{"key": "value"}
	Set(wantLabels)
	if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
		t.Errorf("parent goroutine's profile labels: got %v, want %v", gotLabels, wantLabels)
	}
	go func() {
		if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
			t.Errorf("child goroutine's profile labels: got %v, want %v", gotLabels, wantLabels)
		}
		sync <- struct{}{}
	}()
	<-sync

	wantLabels = map[string]string{}

	Set(wantLabels)
	if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
		t.Errorf("Expected parent goroutine's profile labels to be empty, got %v", gotLabels)
	}
	go func() {
		if gotLabels := Get(); !reflect.DeepEqual(gotLabels, wantLabels) {
			t.Errorf("Expected child goroutine's profile labels to be empty, got %v", gotLabels)
		}
		sync <- struct{}{}
	}()
	<-sync
}
