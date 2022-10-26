echo "NODE \t\t\t\t\t\tCAPACITY \tMAX CAPACITY"

kubectl get nodes -o json | \
	jq -r '.items[].metadata | .name + " " + .annotations["rok.arrikto.com/capacity"] + " " + .labels["rok.arrikto.com/max-instance-capacity"]'\
	| awk -F ' ' '{print $1 "\t"   int($2/1024^3) " Gi\t\t"  int($3/1024^3) " Gi "}'
