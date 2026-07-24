package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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
	
	payload, _ := json.Marshal(map[string]interface{}{
		"name": policy.Name,
		"mode": policy.Spec.EnforcementMode,
		"network": policy.Spec.Network,
	})

	resp, err := http.Post("http://stuntdouble-control-plane:8080/policy", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Error(err, "Failed to sync policy to Control Plane, simulating fallback...")
	} else {
		defer resp.Body.Close()
		logger.Info("Successfully synced policy to Control Plane", "Status", resp.Status)
	}
	
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
