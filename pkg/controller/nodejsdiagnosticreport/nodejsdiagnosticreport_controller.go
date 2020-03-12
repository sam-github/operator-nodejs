package nodejsdiagnosticreport

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/sam-github/operator-nodejs/internal/kwrap"
	opnodejsv1alpha1 "github.com/sam-github/operator-nodejs/pkg/apis/opnodejs/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nodejsdiagnosticreport")

// Add creates a new NodejsDiagnosticReport Controller and adds it to the
// Manager. The Manager will set fields on the Controller and Start it when the
// Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNodejsDiagnosticReport{
		client: mgr.GetClient(),
		config: mgr.GetConfig(),
		scheme: mgr.GetScheme(),
	}
}

// add a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("nodejsdiagnosticreport-controller", mgr,
		controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	watchNamespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Error(err, "Operator requires WATCH_NAMESPACE")
		panic("Operator requires WATCH_NAMESPACE")
	}

	kind := &source.Kind{Type: &opnodejsv1alpha1.NodejsDiagnosticReport{}}
	handler := &handler.EnqueueRequestForObject{}
	updated := predicate.GenerationChangedPredicate{}
	ns := kwrap.MatchNamespacePredicate(watchNamespace)

	log.Info("Watching namespaces", "namespaces", ns)

	// Watch for changes to primary resource NodejsDiagnosticReport
	err = c.Watch(kind, handler, updated, ns)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNodejsDiagnosticReport implements
// reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNodejsDiagnosticReport{}

// ReconcileNodejsDiagnosticReport reconciles a NodejsDiagnosticReport object
type ReconcileNodejsDiagnosticReport struct {
	// This client, initialized using mgr.Client() above, is a split client that
	// reads objects from the cache and writes to the apiserver
	client   client.Client
	config   *rest.Config
	recorder record.EventRecorder
	scheme   *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NodejsDiagnosticReport object
// and makes changes based on the state read and what is in the
// NodejsDiagnosticReport.Spec
//
// Note: The Controller will requeue the Request to be processed again if the
// returned error is non-nil or Result.Requeue is true, otherwise upon
// completion it will remove the work from the queue.
func (r *ReconcileNodejsDiagnosticReport) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues(
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name,
	)
	reqLogger.Info("Reconciling NodejsDiagnosticReport")

	// Fetch the NodejsDiagnosticReport instance
	instance := &opnodejsv1alpha1.NodejsDiagnosticReport{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Instance not found, deleted before reconciliation - don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue
		return reconcile.Result{}, err
	}

	// Do not reconcile if the dump already started - don't requeue
	if instance.Status.Result != "" {
		return reconcile.Result{}, nil
	}

	// Check if this Pod exists and is running
	where := types.NamespacedName{
		Name: instance.Spec.PodName,
		// XXX(sam) is it reasonable to assume that pod is in request's namespace?
		Namespace: request.Namespace,
	}
	pod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), where, pod)

	if err != nil && errors.IsNotFound(err) {
		// Permanent failure - don't requeue
		reqLogger.Error(err, "Cannot generate report, pod not found",
			"Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)

		instance.Status.Result = "PodNotFound"
		r.client.Status().Update(context.TODO(), instance)
		return reconcile.Result{}, nil
	} else if err != nil {
		// Ephemeral failure - requeue
		instance.Status.Result = "RetryOnError"
		return reconcile.Result{}, err
	} else if pod.Status.Phase != corev1.PodRunning {
		// Permananent failure - don't requeue
		reqLogger.Error(err, "Cannot generate report, pod state is not running",
			"Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name,
			"Pod.Status", pod.Status.Phase)

		instance.Status.Result = "PodNotRunning"
		r.client.Status().Update(context.TODO(), instance)
		return reconcile.Result{}, nil
	}

	cmd := []string{"/bin/sh", "-c", "kill -USR2 $(pidof node)"}

	_, err = kwrap.ExecuteCommandInContainer(r.config, where, cmd)
	if err != nil {
		// Permananent failure - don't requeue
		log.Error(err, "Exec report in container", "cmd", cmd)

		instance.Status.Result = "TriggerFailed"
		r.client.Status().Update(context.TODO(), instance)
		return reconcile.Result{}, nil
	}

	// Success - don't requeue
	log.Info("Exec report in container", "cmd", cmd)
	instance.Status.Result = "Complete"
	r.client.Status().Update(context.TODO(), instance)
	return reconcile.Result{}, nil
}
