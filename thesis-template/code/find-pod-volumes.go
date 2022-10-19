func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*v1.PersistentVolumeClaim, node *v1.Node, simulateUnpinnedVolumes bool) (podVolumes *PodVolumes, reasons ConflictReasons, err error) {
	...
	rokClaimsSimulatedUnpinned := []*v1.PersistentVolumeClaim{}
	// Check PV node affinity on bound volumes
	if len(boundClaims) > 0 {
		// Bound PVCs provisioned by rok.arrikto.com that are simulated as unpinned,
		// will skip checkBoundClaims check, and instead will be treated as volumes
		// that need to be provisioned.
		// Thus, we accumulate them in rokClaimsSimulatedUnpinned array,
		// and then append then to claimsToProvision.
		if simulateUnpinnedVolumes {
			// boundClaimsFiltered is the boundClaims after filtering out Rok PVCs simulated as unpinned.
			boundClaimsFiltered := []*v1.PersistentVolumeClaim{}

			for _, claim := range boundClaims {
				if isRok, _ := b.isRokClaim(claim); isRok {
					rokClaimsSimulatedUnpinned = append(rokClaimsSimulatedUnpinned, claim)
					klog.V(4).Infof("Bound claim %q is provisioned by 'rok.arrikto.com' and simulated as unpinned: adding it to the list of volumes that need provisioning", getPVCName(claim))
				} else {
					boundClaimsFiltered = append(boundClaimsFiltered, claim)

				}
			}
			boundClaims = boundClaimsFiltered
		}

		boundVolumesSatisfied, boundPVsFound, err = b.checkBoundClaims(boundClaims, node, podName)
		if err != nil {
			return
		}
	}

	// Find matching volumes and node for unbound claims
	if (len(claimsToBind) > 0) || (len(rokClaimsSimulatedUnpinned) > 0) {
		...
		// Find matching volumes
		if len(claimsToFindMatching) > 0 {
			var unboundClaims []*v1.PersistentVolumeClaim
			unboundVolumesSatisfied, staticBindings, unboundClaims, err = b.findMatchingVolumes(pod, claimsToFindMatching, node)
			if err != nil {
				return
			}
			claimsToProvision = append(claimsToProvision, unboundClaims...)
		}

		// Check for claims to provision. This is the first time where we potentially
		// find out that storage is not sufficient for the node.
		claimsToProvision = append(claimsToProvision, rokClaimsSimulatedUnpinned...)
		if len(claimsToProvision) > 0 {
			unboundVolumesSatisfied, sufficientStorage, dynamicProvisions, err = b.checkVolumeProvisions(pod, claimsToProvision, node)
		}
	}
	return
}