@@ -235,7 +240,7 @@ func NewVolumeBinder(
 // FindPodVolumes finds the matching PVs for PVCs and nodes to provision PVs
 // for the given pod and node. If the node does not fit, confilict reasons are
 // returned.
-func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*v1.PersistentVolumeClaim, node *v1.Node) (podVolumes *PodVolumes, reasons ConflictReasons, err error) {
+func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*v1.PersistentVolumeClaim, node *v1.Node, simulateUnpinnedVolumes bool) (podVolumes *PodVolumes, reasons ConflictReasons, err error) {
 	podVolumes = &PodVolumes{}
 	podName := getPodName(pod)
 
@@ -292,8 +297,31 @@ func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*
 		podVolumes.DynamicProvisions = dynamicProvisions
 	}()
 
+	rokClaimsSimulatedUnpinned := []*v1.PersistentVolumeClaim{}
 	// Check PV node affinity on bound volumes
 	if len(boundClaims) > 0 {
+		// Bound PVCs provisioned by rok.arrikto.com that are simulated as unpinned,
+		// will skip checkBoundClaims check, and instead will be treated as volumes
+		// that need to be provisioned.
+		// Thus, we accumulate them in rokClaimsSimulatedUnpinned array,
+		// and then append then to claimsToProvision.
+		if simulateUnpinnedVolumes {
+			// boundClaimsFiltered is the boundClaims after filtering out Rok PVCs simulated as unpinned.
+			boundClaimsFiltered := []*v1.PersistentVolumeClaim{}
+
+			for _, claim := range boundClaims {
+				if isRok, _ := b.isRokClaim(claim); isRok {
+					rokClaimsSimulatedUnpinned = append(rokClaimsSimulatedUnpinned, claim)
+					klog.V(4).Infof("Bound claim %q is provisioned by 'rok.arrikto.com' and simulated as unpinned: adding it to the list of volumes that need provisioning", getPVCName(claim))
+				} else {
+					boundClaimsFiltered = append(boundClaimsFiltered, claim)
+
+				}
+			}
+
+			boundClaims = boundClaimsFiltered
+		}
+
 		boundVolumesSatisfied, boundPVsFound, err = b.checkBoundClaims(boundClaims, node, podName)
 		if err != nil {
 			return
@@ -301,7 +329,7 @@ func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*
 	}
 
 	// Find matching volumes and node for unbound claims
-	if len(claimsToBind) > 0 {
+	if (len(claimsToBind) > 0) || (len(rokClaimsSimulatedUnpinned) > 0) {
 		var (
 			claimsToFindMatching []*v1.PersistentVolumeClaim
 			claimsToProvision    []*v1.PersistentVolumeClaim
@@ -333,6 +361,7 @@ func (b *volumeBinder) FindPodVolumes(pod *v1.Pod, boundClaims, claimsToBind []*
 
 		// Check for claims to provision. This is the first time where we potentially
 		// find out that storage is not sufficient for the node.
+		claimsToProvision = append(claimsToProvision, rokClaimsSimulatedUnpinned...)
 		if len(claimsToProvision) > 0 {
 			unboundVolumesSatisfied, sufficientStorage, dynamicProvisions, err = b.checkVolumeProvisions(pod, claimsToProvision, node)
 			if err != nil {
@@ -850,6 +879,8 @@ func (b *volumeBinder) findMatchingVolumes(pod *v1.Pod, claimsToBind []*v1.Persi
 func (b *volumeBinder) checkVolumeProvisions(pod *v1.Pod, claimsToProvision []*v1.PersistentVolumeClaim, node *v1.Node) (provisionSatisfied, sufficientStorage bool, dynamicProvisions []*v1.PersistentVolumeClaim, err error) {
 	podName := getPodName(pod)
 	dynamicProvisions = []*v1.PersistentVolumeClaim{}
+	// Rok PVCs
+	rokClaims := []*v1.PersistentVolumeClaim{}
 
 	// We return early with provisionedClaims == nil if a check
 	// fails or we encounter an error.
@@ -876,6 +907,13 @@ func (b *volumeBinder) checkVolumeProvisions(pod *v1.Pod, claimsToProvision []*v
 			return false, true, nil, nil
 		}
 
+		// PVCs provisioned by rok.arrikto.com
+		// Accumulate Rok PVCs in an array and skip hasEnoughCapacity
+		if provisioner == RokProvisioner {
+			rokClaims = append(rokClaims, claim)
+			continue
+		}
+
 		// Check storage capacity.
 		sufficient, err := b.hasEnoughCapacity(provisioner, claim, class, node)
 		if err != nil {
@@ -889,11 +927,43 @@ func (b *volumeBinder) checkVolumeProvisions(pod *v1.Pod, claimsToProvision []*v
 		dynamicProvisions = append(dynamicProvisions, claim)
 
 	}
+
+	// Check Rok storage capacity.
+	sufficient, err := b.hasRokEnoughCapacity(rokClaims, node)
+	if err != nil {
+		return false, false, nil, err
+	}
+	if !sufficient {
+		// hasRokEnoughCapacity logs an explanation.
+		return true, false, nil, nil
+	}
+	// Apend Rok volumes to dynamicProvisions
+	dynamicProvisions = append(dynamicProvisions, rokClaims...)
+
 	klog.V(4).Infof("Provisioning for %d claims of pod %q that has no matching volumes on node %q ...", len(claimsToProvision), podName, node.Name)
 
 	return true, true, dynamicProvisions, nil
 }