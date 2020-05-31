package xdaemonset

import (
	"fmt"
	"sort"
	"strconv"
	//"reflect"
	"context"

	dsv1alpha1 "github.com/wu0407/daemonset-operator/pkg/apis/ds/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"github.com/go-logr/logr"
	"time"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
)

const (
	finalizerName = "xdaemonset.xiaoqing.com"
)

var log = logf.Log.WithName("controller_xdaemonset")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Xdaemonset Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileXdaemonset{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("xdaemonset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	predic := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
		  // Ignore updates to CR status in which case metadata.Generation does not change
		  return e.MetaOld.GetGeneration() != e.MetaNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
		  // Evaluates to false if the object has been confirmed deleted.
		  return !e.DeleteStateUnknown
		},
	}

	// Watch for changes to primary resource Xdaemonset
	err = c.Watch(&source.Kind{Type: &dsv1alpha1.Xdaemonset{}}, &handler.EnqueueRequestForObject{}, predic)
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Xdaemonset
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dsv1alpha1.Xdaemonset{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileXdaemonset implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileXdaemonset{}

// ReconcileXdaemonset reconciles a Xdaemonset object
type ReconcileXdaemonset struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Xdaemonset object and makes changes based on the state read
// and what is in the Xdaemonset.Spec, handle add\update\del event, filter event use Predicate in add().
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileXdaemonset) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Xdaemonset")

	// Fetch the Xdaemonset instance
	instance := &dsv1alpha1.Xdaemonset{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	// Check if the xdaemonset instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isXdaemonSetMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isXdaemonSetMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), finalizerName) {
			// Run finalization logic for finalizerName. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeXdaemonSet(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}

			// Remove finalizerName. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(instance, finalizerName)
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !contains(instance.GetFinalizers(), finalizerName) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	dsList, err := r.getDaemonsetList(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	if dsList == nil {
		return reconcile.Result{}, nil
	}

	leng := len(dsList.Items)
	switch  {
	case leng == 0:
		// Define a new daemonset object
		ds := newDaemonSetForCR(instance)

		// Set Xdaemonset instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, ds, r.scheme); err != nil {
			reqLogger.Error(err, "err in controllerutil.SetControllerReference(instance, ds, r.scheme) 0")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Creating a new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
		err = r.client.Create(context.TODO(), ds)
		if err != nil {
			reqLogger.Error(err, "err in r.client.Create(context.TODO(), ds) 0")
			return reconcile.Result{}, err
		}
		// daemonset created successfully - don't requeue
		return reconcile.Result{}, nil
	case leng == 1:
		// daemonset already exists and spec not change - don't requeue
		if apiequality.Semantic.DeepEqual(dsList.Items[0].Spec, instance.Spec.DaemonSetSpec) {
			fmt.Println("equal")
			return reconcile.Result{}, nil
		}

		//daemonset spec change, create new daemonset, next delete old daemonset 
		newds := newDaemonSetForCR(instance)
		// Set Xdaemonset instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, newds, r.scheme); err != nil {
			reqLogger.Error(err, "err in controllerutil.SetControllerReference")
			return reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new Daemonset", "Daemonset.Namespace", newds.Namespace, "Daemonset.Name", newds.Name)
		err = r.client.Create(context.TODO(), newds)
		if err != nil {
			reqLogger.Error(err, "err in r.client.Create(context.TODO(), newds)")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	default:
		var dss dsslicetype
		dss = dsList.Items
		sort.Sort(dss)
		//new daemonset is ready
		if dss[leng - 1].Status.NumberReady == dss[leng - 1].Status.DesiredNumberScheduled && dss[leng - 1].Status.DesiredNumberScheduled == dss[leng - 1].Status.CurrentNumberScheduled {
			//delete old daemonset
			err = r.client.Delete(context.TODO(), &dss[leng - 2])
			if err != nil {
				reqLogger.Error(err, "err in r.client.Delete(context.TODO(), &dss[leng - 2])")
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		return reconcile.Result{Requeue: true}, nil
	}
	
}

// newDaemonSetForCR returns a DaemonSet with the same name/namespace as the cr
func newDaemonSetForCR(cr *dsv1alpha1.Xdaemonset) *appsv1.DaemonSet {
	labels := make(map[string]string)
	for key, value := range cr.Spec.Template.GetLabels() {
		labels[key] = value
	}
	hash := time.Now().Format("060102150405")
	labels["pod-template-hash"] = hash

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-" + hash,
		//	Name:      cr.Name + "-ds",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: cr.Spec.DaemonSetSpec,
	}
}

func (r *ReconcileXdaemonset) finalizeXdaemonSet(reqLogger logr.Logger, d *dsv1alpha1.Xdaemonset) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	reqLogger.Info("Successfully finalized Xdeamonset")
	return nil
}

func (r *ReconcileXdaemonset) addFinalizer(reqLogger logr.Logger, d *dsv1alpha1.Xdaemonset) error {
	reqLogger.Info("Adding Finalizer for the Xdeamonset")
	controllerutil.AddFinalizer(d, finalizerName)

	// Update CR
	err := r.client.Update(context.TODO(), d)
	if err != nil {
		reqLogger.Error(err, "Failed to update Xdeamonset with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

//getDaemonset return  owner by Xdaemonset  daemonset list
func (r *ReconcileXdaemonset) getDaemonsetList(d *dsv1alpha1.Xdaemonset) (dslist *appsv1.DaemonSetList, err error) {
	daemonsetSelector, err := metav1.LabelSelectorAsSelector(d.Spec.Selector)
	if err != nil {
		return nil, err
	}

	//todo filter by metadata.ownerReferences
	//listOpts := []client.ListOption{
	//	client.InNamespace(d.Namespace),
	//	client.MatchingLabels(d.Spec.Selector.MatchLabels),
	//}
	//err = r.client.List(context.TODO(), dslist, listOpts...)
	dslist = &appsv1.DaemonSetList{}
	err = r.client.List(context.TODO(), dslist, &client.ListOptions{
		Namespace: d.Namespace,
		LabelSelector: daemonsetSelector,
	})
	return dslist, err
}

type dsslicetype []appsv1.DaemonSet

func (a dsslicetype) Len() int           { return len(a) }
func (a dsslicetype) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a dsslicetype) Less(i, j int) bool {
	if a[i].CreationTimestamp.Equal(&a[j].CreationTimestamp) {
		return a[i].Name < a[j].Name
	}
	return a[i].CreationTimestamp.Before(&a[j].CreationTimestamp)
	
	//timei, _ := strconv.Atoi(a[i].CreationTimestamp.Time.Format("20060102150405"))
	//timej, _ := strconv.Atoi(a[j].CreationTimestamp.Time.Format("20060102150405"))
	//return timei < timej
}