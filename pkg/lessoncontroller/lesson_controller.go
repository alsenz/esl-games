/*
Copyright 2021.
*/

package lessoncontroller

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	eslgamesalsenzgithubcomv1alpha1 "github.com/alsenz/esl-games/pkg/k8s/api/v1alpha1"
)

// LessonReconciler reconciles a Lesson object
type LessonReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=esl-games.alsenz.github.com,resources=lessons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=esl-games.alsenz.github.com,resources=lessons/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=esl-games.alsenz.github.com,resources=lessons/finalizers,verbs=update

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
	logger := r.Log.WithValues("lesson", req.NamespacedName)
	logger.Info("Reconciling Lesson Resource")

	instance := &eslgamesalsenzgithubcomv1alpha1.Lesson{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	err = r.ensureLatestPod(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	//TODO we need to set some kind of timer to ensure that old pods (and old lesson resources) are destroyed.
	//TODO TODO

	if instance.Spec.Auth.Login.Required {
		//TODO we need to do something to ensure that we create an oauth proxy here...
		logger.Info("Lesson request with login, this is not implemented yet!")
	}

	return ctrl.Result{}, nil
}

func (r *LessonReconciler) ensureLatestPod(instance *eslgamesalsenzgithubcomv1alpha1.Lesson) error {
	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set Presentation instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return err
	}
	// Check if this Pod already exists
	found := &corev1.Pod{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return err
		}

		// Pod created successfully - don't requeue
		return nil
	} else if err != nil {

		return err
	}

	return nil
}

//TODO newServiceForPod is next step
//TODO newIngressRule too

// newPodForCR returns a deployment with a games-server for the lesson
func newPodForCR(cr *eslgamesalsenzgithubcomv1alpha1.Lesson) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "games-server",
					Image: "games-server-image", //TODO this image tag may very well need to be read from a config map...
					Args: []string{}, //TODO need to add a bunch of arguments that set up the game server...
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *LessonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eslgamesalsenzgithubcomv1alpha1.Lesson{}).
		Complete(r)
}
