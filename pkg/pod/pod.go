package pod

import (
	"fmt"
	"github.com/aidarkhanov/nanoid/v2"
)

type Pod struct {
	Id        string
	client    *containerd.Client
	ctx       *context.context
	container *containerd.Container
}

type RunningPod struct {
	Pod         *Pod
	task        *containerd.Task
	exitStatusC <-chan containerd.exitStatus
}

func NewPod(registryImage string, name string) (*Pod, error) {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return nil, err
	}

	ctx := namespaces.WithNameSpace(context.Background(), "cuba")

	image, err := client.Pull(ctx, registryImage, containerd.WithPullUnpack)

	if err != nil {
		return nil, err
	}

	id := generateNewId(name)

	container, err := client.NewContainer(
		ctx,
		id,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(id+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	if err != nil {
		return nil, err
	}

	return &Pod{
		Id:        id,
		container: &container,
		ctx:       &ctx,
		client:    &client,
	}, nil
}

func generateNewId(name string) string {
	id := uuid.New().String()

	return fmt.Sprintf("%s-%s", name, id)
}

func (pod *Pod) Run() (*RunningPod, error) {
	task, err := (*pod.container).NewTask(*pod.ctx, cio.NewCreator(cio.WithStdio))

	if err != nil {
		return nil, err
	}

	exitStatusC, err := task.Wait(*pod.ctx)
	if err != nil {
		return nil, err
	}

	if err := task.Start(*pod.ctx); err != nil {
		return nil, err
	}

	return &RunningPod{
		Pod:         pod,
		task:        &task,
		exitStatusC: exitStatusC,
	}, nil
}

func (pod *Pod) Kill() (uint32, error) {
	if err := (*pod.task).Kill(*pod.Pod.ctx, syscall.SIGTERM); err != nil {
		return 0, err
	}

	status := <-(*pod.exitStatusC)
	code, _, err := status.Result()

	if err != nil {
		return 0, err
	}

	(*pod.task).Delete(*pod.ctx)

	return code, nil
}

func (pod *Pod) Delete() {
	(*pod.container).Delete(*pod.ctx, containerd.WithSnaptshotCleanup)
	pod.client.Close()
}
