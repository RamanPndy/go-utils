package goutils

import (
	"encoding/json"
)

// --- K8s API list response types ---

type k8sListResponse struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Items      []json.RawMessage `json:"items"`
}

type k8sObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp string            `json:"creationTimestamp"`
}

type k8sObject struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   k8sObjectMeta   `json:"metadata"`
	Spec       json.RawMessage `json:"spec"`
	Status     json.RawMessage `json:"status"`
}

type k8sEvent struct {
	APIVersion     string        `json:"apiVersion"`
	Kind           string        `json:"kind"`
	Metadata       k8sObjectMeta `json:"metadata"`
	Type           string        `json:"type"`
	Reason         string        `json:"reason"`
	Message        string        `json:"message"`
	InvolvedObject struct {
		Kind      string `json:"kind"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		UID       string `json:"uid"`
	} `json:"involvedObject"`
	Source struct {
		Component string `json:"component"`
	} `json:"source"`
	FirstTimestamp string `json:"firstTimestamp"`
	LastTimestamp  string `json:"lastTimestamp"`
	Count          int    `json:"count"`
}

type k8sPodStatus struct {
	Phase             string          `json:"phase"`
	ContainerStatuses json.RawMessage `json:"containerStatuses"`
}

type k8sPod struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Metadata   k8sObjectMeta `json:"metadata"`
	Spec       struct {
		NodeName string `json:"nodeName"`
	} `json:"spec"`
	Status k8sPodStatus `json:"status"`
}

type k8sJob struct {
	Metadata k8sObjectMeta `json:"metadata"`
	Status   struct {
		StartTime      string `json:"startTime"`
		CompletionTime string `json:"completionTime"`
		Succeeded      int    `json:"succeeded"`
		Failed         int    `json:"failed"`
	} `json:"status"`
}
