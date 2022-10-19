/*
 * This file is part of Rok.
 *
 * Copyright © 2022 Arrikto Inc.  All Rights Reserved.
 */

package main

const LoggerName = "rok-scheduler-webhook"
const PodMutatedLabel = "rok-scheduler-webhook.arrikto.com/mutated"

// podHandler handles Pods
type podHandler struct {
	Client           client.Client
	decoder          *admission.Decoder
	optOutAnnotation string
	optInNamespaces  []glob.Glob
	schedulerName    string
}

// podHandler implements admission.DecoderInjector.
// A decoder will be automatically injected.
// InjectDecoder injects the decoder.
func (p *podHandler) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}

func (p *podHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := p.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	// Set Pod namespace to AdmissionRequest namespace as Pod namespace might be empty.
	pod.Namespace = req.Namespace
	podName := pod.Name
	podGenerateName := pod.GenerateName

	if p.skipAdmission(pod) {
		log.Log.WithName(LoggerName).WithValues("namespace", pod.Namespace,
			"name", podName, "generateName", podGenerateName).Info("skipping pod mutation")
		return admission.Allowed("Admission skipped")
	}

	// Mutate the pod
	addSchedulerName(pod, p.schedulerName)
	addLabel(pod, PodMutatedLabel, "true")
	log.Log.WithName(LoggerName).WithValues("namespace", pod.Namespace,
		"name", podName, "generateName", podGenerateName).Info("adding scheduler name to pod", "schedulerName", p.schedulerName)

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	// JSON patches are generated automatically
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// skipAdmission checks if admission should be skipped for the specific Pod.
// This can happen because:
// - The Pod has an opt-out annotation
// - The Pod is not in the opt-in namespaces
func (p *podHandler) skipAdmission(pod *corev1.Pod) bool {
	podName := pod.Name
	podGenerateName := pod.GenerateName
	if hasAnnotationKey(pod, p.optOutAnnotation) {
		log.Log.WithName(LoggerName).WithValues("namespace", pod.Namespace, "name", podName,
			"generateName", podGenerateName).Info("pod has skip mutation annotation key", "annotation", p.optOutAnnotation)
		return true

	}
	if matchesAnyGlob(p.optInNamespaces, pod.Namespace) {
		return false

	}
	log.Log.WithName(LoggerName).WithValues("namespace", pod.Namespace,
		"name", podName, "generateName", podGenerateName).Info("pod's namespace didn't match any of the namespaces the webhook mutates")
	return true
}
