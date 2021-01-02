/*
Copyright 2021.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	eslgamesalsenzgithubcomv1alpha1 "github.com/alsenz/esl-games/api/v1alpha1"
)

// LessonReconciler reconciles a Lesson object
type LessonReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=esl-games.alsenz.github.com.esl-games.alsenz.github.com,resources=lessons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=esl-games.alsenz.github.com.esl-games.alsenz.github.com,resources=lessons/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=esl-games.alsenz.github.com.esl-games.alsenz.github.com,resources=lessons/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Lesson object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *LessonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("lesson", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LessonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eslgamesalsenzgithubcomv1alpha1.Lesson{}).
		Complete(r)
}
