func (a *StaticAutoscaler) obtainNodeLists(cp cloudprovider.CloudProvider) ([]*apiv1.Node, []*apiv1.Node, errors.AutoscalerError) {
	allNodes, err := a.AllNodeLister().List()
	readyNodes, err := a.ReadyNodeLister().List()
        ...
	allNodes, readyNodes = gpu.FilterOutNodesWithUnreadyGpus(cp.GPULabel(), allNodes, readyNodes)
	allNodes, readyNodes = taints.FilterOutNodesWithIgnoredTaints(a.ignoredTaints, allNodes, readyNodes)
	allNodes, readyNodes = FilterOutNodesWithUnreadyCSI(allNodes, readyNodes)
	return allNodes, readyNodes, nil
}
