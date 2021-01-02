/*
Copyright 2021.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type LoginSpec struct {
	Required bool   `json:"required,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type AuthSpec struct {
	Passcode     string    `json:"passcode,omitempty"`
	Login        LoginSpec `json:"login,omitempty"`
	FallbackMode bool      `json:"fallbackMode,omitempty"`
	Public       bool      `json:"public,omitempty"`
}

// LessonSpec defines the desired state of Lesson
type LessonSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Lesson. Edit Lesson_types.go to remove/update
	Code      string   `json:"code,omitempty"`
	LeaderUri string   `json:"leaderUri,omitempty"`
	PlanUri   string   `json:"planUri,omitempty"`
	Auth      AuthSpec `json:"auth,omitempty"`
}

// LessonStatus defines the observed state of Lesson
type LessonStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Lesson is the Schema for the lessons API
type Lesson struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LessonSpec   `json:"spec,omitempty"`
	Status LessonStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LessonList contains a list of Lesson
type LessonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lesson `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lesson{}, &LessonList{})
}
