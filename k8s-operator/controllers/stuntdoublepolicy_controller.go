package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stuntdoublev1alpha1 "github.com/itsrohan-lang/stuntdouble/k8s-operator/api/v1alpha1"
)

// StuntDoublePolicyReconciler reconciles a StuntDoublePolicy object
type StuntDoublePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stuntdouble.io,resources=stuntdoublepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stuntdouble.io,resources=stuntdoublepolicies/status,verbs=get;update;patch

func (r *StuntDoublePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var policy stuntdoublev1alpha1.StuntDoublePolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Syncing StuntDoublePolicy to global Control Plane", "PolicyName", policy.Name)
	// Integration with StuntDouble Control Plane happens here.
	
	policy.Status.ActiveAgents = 1
	if err := r.Status().Update(ctx, &policy); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *StuntDoublePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stuntdoublev1alpha1.StuntDoublePolicy{}).
		Complete(r)
}
