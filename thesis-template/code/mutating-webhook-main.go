func main() {
	...
	// Get command line parameters
	var port int
	var tlsDir, certName, keyName, optOutAnnotation, optInNamespaces, schedulerName string

	flag.IntVar(&port, "port", 443, "Webhook server port.")
	flag.StringVar(&tlsDir, "tls-dir", "/etc/webhook/certs", "Folder containing the X509 certs for the webhook.")
	flag.StringVar(&certName, "cert-name", "cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&keyName, "key-name", "key.pem", "File containing the x509 private key.")
	flag.StringVar(&optOutAnnotation, "annotation-optout", "arrikto.com/skip-rok-scheduler-webhook", "Annotation key that if present, the pod will not get mutated.")
	// Default namespaces-optin: "*" allows mutation in all namespaces
	flag.StringVar(&optInNamespaces, "namespaces-optin", "*", "A comma separated list of namespaces. If a pod matches against these namespaces it will get mutated. Globs can be provided.")
	flag.StringVar(&schedulerName, "scheduler-name", "rok-scheduler", "The name of the scheduler that will be set on the pod.")
	flag.Parse()
	...
	// Compile globs from the namespaces-optin comma separated list
	optInNamespacesGlobs, err := compileGlobsFromCommaSeparatedList(optInNamespaces)
	if err != nil {
		logger.Error(err, "failed to parse the namespace globs list")
		os.Exit(1)
	}

	// Setup webhooks
	logger.Info("setting up webhook server") hookServer := mgr.GetWebhookServer()
	hookServer.Port = port
	hookServer.CertDir = tlsDir
	hookServer.CertName = certName
	hookServer.KeyName = keyName

	logger.Info("registering webhook to the webhook server")
	handler := &podHandler{
		Client:           mgr.GetClient(),
		optOutAnnotation: optOutAnnotation,
		optInNamespaces:  optInNamespacesGlobs,
		schedulerName:    schedulerName,
	}
	hookServer.Register("/mutate", &webhook.Admission{Handler: handler})

	logger.Info("starting manager")
	...
}