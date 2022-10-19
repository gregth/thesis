const (
	// RokCapacityAnnotation is the key of the annotation that exposes the free capacity of Rok on the node
	RokCapacityAnnotation = "rok.arrikto.com/capacity"
	// RokInstanceCapacityLabel is the key of the label that exposes the max capacity of Rok on the node
	RokInstanceCapacityLabel = "rok.arrikto.com/max-instance-capacity"
	// NodeMaxCapacity is the maximum possible storage capacity a node can provide.
	NodeMaxCapacity = math.MaxInt64
)
 
// SanitizeRokStorageAnnotations sets appropriate storage annotations for Rok
func sanitizeRokStorageAnnotations(labels map[string]string, annotations map[string]string, nodeName string) map[string]string {
	capacity := strconv.Itoa(NodeMaxCapacity)
	if val, found := labels[RokInstanceCapacityLabel]; found {
		klog.V(6).Infof("Found label '%s: %q' on template node %q", RokInstanceCapacityLabel, val, nodeName)
		capacity = val
	} else {
		klog.V(6).Infof("Label %q not found on template node %q: Setting capacity to maximum possible node capacity", RokInstanceCapacityLabel, nodeName)
	}
	klog.V(6).Infof("Setting annotation '%s: %q' on template node %q", RokCapacityAnnotation, capacity, nodeName)

	// We expect the annotations to be nil in case the template is created from a cloud provider plugin
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[RokCapacityAnnotation] = capacity
	return annotations
}