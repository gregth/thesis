// hasRokEnoughCapacity checks whether Rok has enough capacity left for all the
// provided volumes on the given node.
func (b *volumeBinder) hasRokEnoughCapacity(claims []*v1.PersistentVolumeClaim, node *v1.Node) (bool, error) {
	// Rok specific capacity tracking
	claimNames := []string{}

	// Sum the requested capacities
	totalRequestedCapacity := int64(0)
	for _, claim := range claims {
		pvcName := getPVCName(claim)
		quantity, ok := claim.Spec.Resources.Requests[v1.ResourceStorage]
		if !ok {
			// If !ok no capacity requested
			klog.V(4).Infof("PVC %q has no capacity request", pvcName)
			continue
		}
		sizeInBytes := quantity.Value()
		klog.V(4).Infof("PVC %q capacity request: %d bytes", pvcName, sizeInBytes)
		claimNames = append(claimNames, pvcName)
		// Check for overflow
		if (totalRequestedCapacity + sizeInBytes) < totalRequestedCapacity {
			klog.V(4).Infof("Overflow while calculating total requested capacity for Rok PVCs %q", strings.Join(claimNames, ", "))
			return false, fmt.Errorf("integer overflow")
		}
		totalRequestedCapacity += sizeInBytes
	}

	// Get the free space from the node API object
	nodeCapacity, ok := node.ObjectMeta.Annotations[RokCapacityAnnotation]
	if !ok {
		// No annotation found, treat this as no available capacity
		klog.V(4).Infof("Node %q has no '%s' annotation: Assuming Rok is not available on this node", node.ObjectMeta.Name, RokCapacityAnnotation)
		return false, nil
	}
	nodeCapacityInBytes, err := strconv.ParseInt(nodeCapacity, 10, 64)
	if err != nil {
		klog.V(4).Infof("Error while converting capacity string %q to bytes: invalid format", nodeCapacity)
		return false, err
	}

	klog.V(4).Infof("Rok PVCs %q total capacity request: %d bytes, free storage capacity on node %q: %d bytes", strings.Join(claimNames, ", "), totalRequestedCapacity, node.ObjectMeta.Name, nodeCapacityInBytes)
	if nodeCapacityInBytes >= totalRequestedCapacity {
		// Enough capacity found.
		klog.V(4).Infof("Sufficient free storage capacity for PVCs %q on node %q", strings.Join(claimNames, ", "), node.ObjectMeta.Name)
		return true, nil
	}
	klog.V(4).Infof("Insufficient free storage capacity for PVCs %q on node %q", strings.Join(claimNames, ", "), node.ObjectMeta.Name)
	return false, nil
}