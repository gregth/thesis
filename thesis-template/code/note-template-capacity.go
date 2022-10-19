 func sanitizeTemplateNode(node *apiv1.Node, nodeGroup string, ignoredTaints taints.TaintKeySet) (*apiv1.Node, errors.AutoscalerError) {
	newNode := node.DeepCopy()
	nodeName := fmt.Sprintf("template-node-for-%s-%d", nodeGroup, rand.Int63())
	...
	newNode.Name = nodeName
	newNode.Spec.Taints = taints.SanitizeTaints(newNode.Spec.Taints, ignoredTaints)
	newNode.ObjectMeta.Annotations = sanitizeRokStorageAnnotations(newNode.Labels, newNode.ObjectMeta.Annotations, newNode.Name)
	return newNode, nil
}
