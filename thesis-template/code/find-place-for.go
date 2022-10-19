func findPlaceFor(removedNode string, pods []*apiv1.Pod, nodes map[string]bool,
	clusterSnapshot ClusterSnapshot, predicateChecker PredicateChecker, oldHints map[string]string, newHints map[string]string, usageTracker *UsageTracker,
	timestamp time.Time) error {
	...
	for _, podptr := range pods {
		...
		klog.V(5).Infof("Looking for place for %s/%s", pod.Namespace, pod.Name)

		// Skip rok-csi-guard
		csiGuardSelector, _ := labels.Parse("app=rok-csi-guard")
		if csiGuardSelector.Matches(labels.Set(pod.Labels)) {
			klog.V(2).Infof("Skipping findPlaceFor for %s pod", pod.Name)
			continue
		}
		...
	return nil
}

func checkPdbs(pods []*apiv1.Pod, pdbs []*policyv1.PodDisruptionBudget) (*drain.BlockingPod, error) {
	for _, pdb := range pdbs {
		// Ignore PDBs for rok-csi-guard pods
		csiGuardSelector, err := labels.Parse("app=rok-csi-guard")
		if err != nil {
			return nil, err
		}
		if csiGuardSelector.Matches(labels.Set(pdb.Labels)) {
			klog.V(2).Infof("Ignoring pod disruption budget  %s/%s", pdb.Namespace, pdb.Name)
			continue
		}
		...
}
