package kwrap

import (
	"bytes"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Where = types.NamespacedName

// ExecuteCommandInContainer will execute a comand inside a pod container
func ExecuteCommandInContainer(config *rest.Config, where Where, cmd []string) (string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("Exec create new client: %w", err)
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(where.Name).
		Namespace(where.Namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Command: cmd,
		Stdout:  true,
		Stderr:  true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("Exec creating SPDY executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		err = fmt.Errorf("Exec cmd %v stderr %v: %w", cmd, stderr.String(), err)
		return stderr.String(), err
	}

	return stdout.String(), nil
}
