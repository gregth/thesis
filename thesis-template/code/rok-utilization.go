// Based on CalculateUtilization(), see simulator/cluster.go
func CalculateUtilizationOfRokStorage(node *apiv1.Node) (utilInfo UtilizationInfo, err error) {
	capacity, foundCapacity := node.ObjectMeta.Annotations[rok.RokCapacityAnnotation]
	maxCapacity, foundMaxCapacity := node.Labels[rok.RokInstanceCapacityLabel]
	if !foundCapacity || !foundMaxCapacity {
		klog.V(3).Infof("Rok doesn't run on node %s", node.Name)
		return UtilizationInfo{RokStorageUtil: 0, ResourceName: rok.RokStorageResource, Utilization: 0}, nil
	}

	// TODO: Restructure this code, used in many places
	nodeCapacityInBytes, err := strconv.ParseInt(capacity, 10, 64)
	if err != nil {
		klog.V(4).Infof("Error while converting capacity string %q to bytes: invalid format", capacity)
		return UtilizationInfo{}, err
	}
	// TODO: Restructure this code, used in many places
	nodeMaxCapacityInBytes, err := strconv.ParseInt(maxCapacity, 10, 64)
	if err != nil {
		klog.V(4).Infof("Error while converting max capacity string %q to bytes: invalid format", maxCapacity)
		return UtilizationInfo{}, err
	}

	util := (float64(nodeMaxCapacityInBytes) - float64(nodeCapacityInBytes)) / float64(nodeMaxCapacityInBytes)
	return UtilizationInfo{RokStorageUtil: util, ResourceName: rok.RokStorageResource, Utilization: util}, nil
}