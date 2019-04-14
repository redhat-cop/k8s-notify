package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	notifyv1 "github.com/redhat-cop/events-notifier/pkg/apis/notify/v1"
	"github.com/redhat-cop/events-notifier/pkg/util"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const finalizer = "finalizers.event.redhat-cop.io"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Notifier Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, sr *util.SharedResources) error {
	return add(mgr, newReconciler(mgr, sr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, sr *util.SharedResources) reconcile.Reconciler {
	return &ReconcileNotifier{client: mgr.GetClient(), scheme: mgr.GetScheme(), sr: sr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("notifier-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Notifier
	err = c.Watch(&source.Kind{Type: &notifyv1.Notifier{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Notifier
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &notifyv1.Notifier{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNotifier{}

// ReconcileNotifier reconciles a Notifier object
type ReconcileNotifier struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	sr     *util.SharedResources
}

// Reconcile reads that state of the cluster for a Notifier object and makes changes based on the state read
// and what is in the Notifier.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNotifier) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithFields(log.Fields{"Request.Namespace": request.Namespace, "Request.Name": request.Name})
	reqLogger.Info("Reconciling Notifier")

	// Fetch the Notifier instance
	instance := notifyv1.Notifier{}
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

	// create a unique identifier for the subscription using `namespace/name`
	id := fmt.Sprintf("%s/%s", instance.Namespace, instance.Name)

	// Log the object for debugging purposes
	out, err := json.Marshal(&instance)
	if err != nil {
		reqLogger.Error(err, "Failed to unmarshall Notifier")
		return reconcile.Result{}, err
	}
	reqLogger.Debug(fmt.Sprintf("Processing Notifier: %s", out))

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// CR is being deleted;

		// remove our finalizer from the list and update it.
		instance.ObjectMeta.Finalizers = util.RemoveString(instance.ObjectMeta.Finalizers, finalizer)
		if err := r.client.Update(context.Background(), &instance); err != nil {
			return reconcile.Result{}, err
		}

		// Unregister notifier
		delete(r.sr.Subscriptions, id)
		reqLogger.Info(fmt.Sprintf("Removed notifier: %s", instance.Name))

		return reconcile.Result{}, nil
	}

	// Ensure finalizer is set
	if !util.ContainsString(instance.ObjectMeta.Finalizers, finalizer) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizer)
		if err := r.client.Update(context.Background(), &instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Register subscription
	r.sr.Notifiers[id] = instance

	return reconcile.Result{}, nil
}
