const (
    // RokCapacityAnnotation is the key of the annotation that exposes the free capacity of Rok on the node
    RokCapacityAnnotation = "rok.arrikto.com/capacity"
 )

func FilterOutNodesWithUnreadyCSI(allNodes, readyNodes []*apiv1.Node) ([]*apiv1.Node, []*apiv1.Node) {
	newAllNodes := make([]*apiv1.Node, 0)
	newReadyNodes := make([]*apiv1.Node, 0)
	nodesWithUnreadyCSI := make(map[string]*apiv1.Node)
	for _, node := range readyNodes {
		hasCapacityAnnotation := false
		if node.Annotations != nil {
			_, hasCapacityAnnotation = node.Annotations[RokCapacityAnnotation]
		}

		if !hasCapacityAnnotation {
			klog.V(3).Infof("Overriding status of node %v, which seems to have unready CSI", node.Name)
			nodesWithUnreadyCSI[node.Name] = kubernetes.GetUnreadyNodeCopy(node)
		} else {
			newReadyNodes = append(newReadyNodes, node)
		}
	}

	// Override any node with unready CSI with its "unready" copy
	for _, node := range allNodes {
		if newNode, found := nodesWithUnreadyCSI[node.Name]; found {
			newAllNodes = append(newAllNodes, newNode)
		} else {
			newAllNodes = append(newAllNodes, node)
		}
	}
	return newAllNodes, newReadyNodes
}