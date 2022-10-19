func (b *volumeBinder) checkVolumeProvisions(pod *v1.Pod, claimsToProvision []*v1.PersistentVolumeClaim, node *v1.Node) (provisionSatisfied, sufficientStorage bool, dynamicProvisions []*v1.PersistentVolumeClaim, err error) {
	podName := getPodName(pod)
	dynamicProvisions = []*v1.PersistentVolumeClaim{}
	// Rok PVCs
	rokClaims := []*v1.PersistentVolumeClaim{}

	for _, claim := range claimsToProvision {
                ...
		class, err := b.classLister.Get(className)
                ...
		provisioner := class.Provisioner
                ...
		// PVCs provisioned by rok.arrikto.com
		// Accumulate Rok PVCs in an array and skip hasEnoughCapacity
		if provisioner == RokProvisioner {
			klog.V(4).Infof("Unbound claim %q is provisioned by 'rok.arrikto.com'", pvcName)
			rokClaims = append(rokClaims, claim)
			continue
		}

		// Check storage capacity.
		sufficient, err := b.hasEnoughCapacity(provisioner, claim, class, node)
                ...
		if !sufficient {
			// hasEnoughCapacity logs an explanation.
			return true, false, nil, nil
		}

		dynamicProvisions = append(dynamicProvisions, claim)
	}

	// Check Rok storage capacity.
	sufficient, err := b.hasRokEnoughCapacity(rokClaims, node)
	if err != nil {
		return false, false, nil, err
	}
	if !sufficient {
		// hasRokEnoughCapacity logs an explanation.
		return true, false, nil, nil
	}

	// Apend Rok volumes to dynamicProvisions
	dynamicProvisions = append(dynamicProvisions, rokClaims...)
	klog.V(4).Infof("Provisioning for %d claims of pod %q that has no matching volumes on node %q ...", len(claimsToProvision), podName, node.Name)
	return true, true, dynamicProvisions, nil
}

const (
	// RokCapacityAnnotation is the key of the annotation that exposes the free capacity of Rok on the node
	RokCapacityAnnotation = "rok.arrikto.com/capacity"
	// RokProvisioner is the name of the Rok provisioner
	RokProvisioner = "rok.arrikto.com"
)
