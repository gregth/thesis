func (sd *ScaleDown) checkNodeUtilization(timestamp time.Time, node *apiv1.Node, nodeInfo *schedulerframework.NodeInfo) (simulator.UnremovableReason, *simulator.UtilizationInfo) {
        ...
	if !sd.isNodeBelowUtilizationThreshold(node, utilInfo) {
		klog.V(4).Infof("Node %s is not suitable for removal - %s utilization too big (%f)", node.Name, utilInfo.ResourceName, utilInfo.Utilization)
		return simulator.NotUnderutilized, &utilInfo
	}

	klog.V(4).Infof("Node %s - %s utilization %f", node.Name, utilInfo.ResourceName, utilInfo.Utilization)

	// Get Rok storage utilization
	if rok.NodeHasRokStorageCapacity(node) {
		utilInfo, err := simulator.CalculateUtilizationOfRokStorage(node)
		if err != nil {
			klog.Warningf("Failed to calculate Rok storage utilization for %s: %v", node.Name, err)
		}
		klog.V(4).Infof("Node %s - %s utilization %f", node.Name, utilInfo.ResourceName, utilInfo.Utilization)

		if !sd.isNodeRokStorageBelowUtilizationThreshold(utilInfo) {
			klog.V(4).Infof("Node %s is not suitable for removal - %s utilization too big (%f)", node.Name, utilInfo.ResourceName, utilInfo.Utilization)
			return simulator.NotUnderutilized, &utilInfo
		}
	}
	return simulator.NoReason, &utilInfo
}