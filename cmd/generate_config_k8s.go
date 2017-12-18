package cmd

import (
	"github.com/spf13/cobra"
	"text/template"
	"os"
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Show a sample k8s/fluentd configuration",
	PreRun: Authorize,
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, _ := template.New("k8s").Parse(k8sTemplateStr)
		f := configTemplateData{AuthHeader: AuthHeader,NovaURL: NovaURL}
		tmpl.Execute(os.Stdout, f)
	},
}

func init() {
	generateConfigCmd.AddCommand(k8sCmd)
}

var k8sTemplateStr = `
# More instructions on how to turn on fluentd integration with k8s
# can be found at https://github.com/splunknova/fluentd

apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: kube-system
  labels:
    k8s-app: fluentd-logging
    version: v1
    kubernetes.io/cluster-service: "true"
spec:
  template:
    metadata:
      labels:
        k8s-app: fluentd-logging
        version: v1
        kubernetes.io/cluster-service: "true"
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: fluentd
        image: polarishq/fluentd_splunknova
        env:
          - name:  SPLUNK_URL
            value: '{{.NovaURL}}:443'
          - name:  SPLUNK_TOKEN
            value: "{{.AuthHeader}}"
          - name: SPLUNK_FORMAT
            value: "nova"
          - name: SPLUNK_URL_PATH
            value: "/v1/events"
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers

`