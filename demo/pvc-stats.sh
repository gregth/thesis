echo "=================\nNAMESPACE: before\n=================\n";
kubectl get pvc -n before;

echo "\n================\nNAMESPACE: after\n================\n";
kubectl get pvc -n after

