package event

import (
	"context"
	"encoding/json"
	"fmt"

	eventv1 "github.com/redhat-cop/events-notifier/pkg/apis/event/v1"
	notifyv1 "github.com/redhat-cop/events-notifier/pkg/apis/notify/v1"
	"github.com/redhat-cop/events-notifier/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_event")

// Add creates a new Service Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, sr *util.SharedResources) error {
	return add(mgr, newReconciler(mgr, sr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, sr *util.SharedResources) reconcile.Reconciler {
	return &ReconcileEvent{client: mgr.GetClient(), scheme: mgr.GetScheme(), sr: sr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("event-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Service
	err = c.Watch(&source.Kind{Type: &corev1.Event{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileEvent{}

// ReconcileRoute reconciles an Event object
type ReconcileEvent struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	sr     *util.SharedResources
}

// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileEvent) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// Fetch the Route svc
	instance := corev1.Event{}
	err := r.client.Get(context.TODO(), request.NamespacedName, &instance)
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

	subscr := r.subscribedTo(&instance)
	if !subscr.Equal(&eventv1.EventSubscription{}) {
		reqLogger.Info(fmt.Sprintf("Notifying of subscribed event: %s", instance.Message))

		// Get notifier associated with subscription
		notifier := &notifyv1.Notifier{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: request.Namespace, Name: subscr.Spec.Notifier}, notifier)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("Notifier `%s` does not exist", subscr.Spec.Notifier))
				return reconcile.Result{}, nil
			}
			// Error reading the object - requeue the request.
			log.Error(err, fmt.Sprintf("Error retreiving Notifier `%s`", subscr.Spec.Notifier))
			return reconcile.Result{}, err
		}

		// Send notification
		err = notifier.GetEventNotifier().Send(instance.Message)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	reqLogger.Info(fmt.Sprintf("Ignoring event %s as there is no matching subscription", instance.ObjectMeta.Name))
	return reconcile.Result{}, nil
}

func (r *ReconcileEvent) subscribedTo(e *corev1.Event) eventv1.EventSubscription {
	reqLogger := log.WithValues()
	var subscribed bool
	var err error
	//var out []byte
	for _, b := range *r.sr.Subscriptions {
		_, err = json.Marshal(b)
		if err != nil {
			reqLogger.Error(err, "Failed to unmarshall EventSubscription")
		}
		// Check for match
		subscribed, err = b.Subscribed(e)
		if err != nil {
			reqLogger.Error(err, "Failed checking subscription")
		}
		if subscribed {
			return b
		}
	}
	return eventv1.EventSubscription{}
}
