package tests

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/golang/glog"
	"github.com/pingcap/tidb-operator/pkg/tkctl/util"
	"github.com/pingcap/tidb-operator/tests/slack"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	DrainerReplicas          int32 = 1
	RunReparoCommandTemplate       = `kubectl exec -it {{ .PodName }} -- bash -c \
"wget http://download.pingcap.org/tidb-binlog-cluster-latest-linux-amd64.tar.gz && \
tar -xzf tidb-binlog-cluster-latest-linux-amd64.tar.gz && \
cd tidb-binlog-cluster-latest-linux-amd64/bin && \
echo '{{ .ReparoConfig }}' > reparo.toml && \
./reparo -config reparo.toml"
`
)

type BackupTarget struct {
	IncrementalType DbType
	TargetCluster   *TidbClusterConfig
	IsAdditional    bool
}

func (t *BackupTarget) GetDrainerConfig(source *TidbClusterConfig, ts string) *DrainerConfig {
	drainerConfig := &DrainerConfig{
		DrainerName:       fmt.Sprintf("to-%s-%s", t.TargetCluster.ClusterName, t.IncrementalType),
		InitialCommitTs:   ts,
		OperatorTag:       source.OperatorTag,
		SourceClusterName: source.ClusterName,
		Namespace:         source.Namespace,
		DbType:            t.IncrementalType,
	}
	if t.IncrementalType == DbTypeMySQL || t.IncrementalType == DbTypeTiDB {
		drainerConfig.Host = fmt.Sprintf("%s.%s.svc.cluster.local",
			t.TargetCluster.ClusterName, t.TargetCluster.Namespace)
		drainerConfig.Port = "4000"
	}
	return drainerConfig
}

func (oa *operatorActions) BackupAndRestoreToMultipleClusters(source *TidbClusterConfig, targets []BackupTarget) error {
	err := oa.DeployAndCheckPump(source)
	if err != nil {
		return err
	}

	err = oa.DeployAdHocBackup(source)
	if err != nil {
		glog.Errorf("cluster:[%s] deploy happen error: %v", source.ClusterName, err)
		return err
	}

	ts, err := oa.CheckAdHocBackup(source)
	if err != nil {
		glog.Errorf("cluster:[%s] deploy happen error: %v", source.ClusterName, err)
		return err
	}

	// Restore can only be done serially due to name collision
	for i := range targets {
		err = oa.CheckTidbClusterStatus(targets[i].TargetCluster)
		if err != nil {
			glog.Errorf("cluster:[%s] deploy faild error: %v", targets[i].TargetCluster.ClusterName, err)
			return err
		}

		err = oa.Restore(source, targets[i].TargetCluster)
		if err != nil {
			glog.Errorf("from cluster:[%s] to cluster [%s] restore happen error: %v",
				source.ClusterName, targets[i].TargetCluster.ClusterName, err)
			return err
		}

		err = oa.CleanRestoreJob(source)
		if err != nil && !releaseIsNotFound(err) {
			glog.Errorf("clean the from cluster:[%s] to cluster [%s] restore job happen error: %v",
				source.ClusterName, targets[i].TargetCluster.ClusterName, err)
		}
	}

	prepareIncremental := func(source *TidbClusterConfig, target BackupTarget) error {
		err = oa.CheckRestore(source, target.TargetCluster)
		if err != nil {
			glog.Errorf("from cluster:[%s] to cluster [%s] restore failed error: %v",
				source.ClusterName, target.TargetCluster.ClusterName, err)
			return err
		}

		if target.IsAdditional {
			// Deploy an additional drainer release
			drainerConfig := target.GetDrainerConfig(source, ts)
			if err := oa.DeployDrainer(drainerConfig, source); err != nil {
				return nil
			}
			if err := oa.CheckDrainer(drainerConfig, source); err != nil {
				return nil
			}
		} else {
			// Enable drainer of the source TiDB cluster release
			if err := oa.DeployAndCheckIncrementalBackup(source, target.TargetCluster, ts); err != nil {
				return err
			}
		}
		return nil
	}

	checkIncremental := func(source *TidbClusterConfig, target BackupTarget) error {
		if target.IncrementalType == DbTypeFile {
			if err := oa.RestoreIncrementalFiles(target.GetDrainerConfig(source, ts), target.TargetCluster); err != nil {
				return err
			}
		}

		if err := oa.CheckDataConsistency(source, target.TargetCluster); err != nil {
			return err
		}
		return nil
	}

	var eg errgroup.Group
	for i := range targets {
		eg.Go(func() error {
			return prepareIncremental(source, targets[i])
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	glog.Infof("waiting 1 minutes to insert into more records")
	time.Sleep(1 * time.Minute)

	glog.Infof("cluster[%s] stop insert data", source.ClusterName)
	oa.StopInsertDataTo(source)

	for i := range targets {
		eg.Go(func() error {
			return checkIncremental(source, targets[i])
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	go oa.BeginInsertDataToOrDie(source)
	err = oa.DeployScheduledBackup(source)
	if err != nil {
		glog.Errorf("cluster:[%s] scheduler happen error: %v", source.ClusterName, err)
		return err
	}

	return oa.CheckScheduledBackup(source)
}

func (oa *operatorActions) BackupAndRestoreToMultipleClustersOrDie(source *TidbClusterConfig, targets []BackupTarget) {
	if err := oa.BackupAndRestoreToMultipleClusters(source, targets); err != nil {
		slack.NotifyAndPanic(err)
	}
}

func (oa *operatorActions) BackupRestore(from, to *TidbClusterConfig) error {

	return oa.BackupAndRestoreToMultipleClusters(from, []BackupTarget{
		{
			TargetCluster:   to,
			IncrementalType: DbTypeTiDB,
			IsAdditional:    false,
		},
	})
}

func (oa *operatorActions) BackupRestoreOrDie(from, to *TidbClusterConfig) {
	if err := oa.BackupRestore(from, to); err != nil {
		slack.NotifyAndPanic(err)
	}
}

func (oa *operatorActions) DeployAndCheckPump(tc *TidbClusterConfig) error {
	if err := oa.DeployIncrementalBackup(tc, nil, false, ""); err != nil {
		return err
	}

	if err := oa.CheckIncrementalBackup(tc, false); err != nil {
		return err
	}
	return nil
}

func (oa *operatorActions) DeployAndCheckIncrementalBackup(from, to *TidbClusterConfig, ts string) error {
	if err := oa.DeployIncrementalBackup(from, to, true, ts); err != nil {
		return err
	}

	if err := oa.CheckIncrementalBackup(from, true); err != nil {
		return err
	}
	return nil
}

func (oa *operatorActions) CheckDataConsistency(from, to *TidbClusterConfig) error {
	fn := func() (bool, error) {
		b, err := to.DataIsTheSameAs(from)
		if err != nil {
			glog.Error(err)
			return false, nil
		}
		if b {
			return true, nil
		}
		return false, nil
	}

	if err := wait.Poll(DefaultPollInterval, 30*time.Minute, fn); err != nil {
		return err
	}
	return nil
}

func (oa *operatorActions) DeployDrainer(info *DrainerConfig, source *TidbClusterConfig) error {
	oa.EmitEvent(source, "DeployDrainer")
	glog.Infof("begin to deploy drainer [%s] namespace[%s], source cluster [%s]", info.DrainerName,
		source.Namespace, source.ClusterName)

	valuesPath, err := info.BuildSubValues(oa.drainerChartPath(source.OperatorTag))
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("helm install %s  --name %s --namespace %s --set-string %s -f %s",
		oa.drainerChartPath(source.OperatorTag), info.DrainerName, source.Namespace, info.DrainerHelmString(nil, source), valuesPath)
	glog.Info(cmd)

	if res, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to deploy drainer [%s/%s], %v, %s",
			source.Namespace, info.DrainerName, err, string(res))
	}

	return nil
}

func (oa *operatorActions) DeployDrainerOrDie(info *DrainerConfig, source *TidbClusterConfig) {
	if err := oa.DeployDrainer(info, source); err != nil {
		slack.NotifyAndPanic(err)
	}
}

func (oa *operatorActions) CheckDrainer(info *DrainerConfig, source *TidbClusterConfig) error {
	glog.Infof("checking drainer [%s/%s]", info.DrainerName, source.Namespace)

	ns := source.Namespace
	stsName := fmt.Sprintf("%s-%s-drainer", source.ClusterName, info.DrainerName)
	fn := func() (bool, error) {
		sts, err := oa.kubeCli.AppsV1().StatefulSets(source.Namespace).Get(stsName, v1.GetOptions{})
		if err != nil {
			glog.Errorf("failed to get drainer StatefulSet %s ,%v", sts, err)
			return false, nil
		}
		if *sts.Spec.Replicas != DrainerReplicas {
			glog.Infof("StatefulSet: %s/%s .spec.Replicas(%d) != %d",
				ns, sts.Name, *sts.Spec.Replicas, DrainerReplicas)
			return false, nil
		}
		if sts.Status.ReadyReplicas != DrainerReplicas {
			glog.Infof("StatefulSet: %s/%s .state.ReadyReplicas(%d) != %d",
				ns, sts.Name, sts.Status.ReadyReplicas, DrainerReplicas)
		}
		return true, nil
	}

	err := wait.Poll(DefaultPollInterval, DefaultPollTimeout, fn)
	if err != nil {
		return fmt.Errorf("failed to install drainer [%s/%s], %v", source.Namespace, info.DrainerName)
	}

	return nil
}

func (oa *operatorActions) RestoreIncrementalFiles(from *DrainerConfig, to *TidbClusterConfig) error {
	glog.Infof("restoring incremental data from drainer [%s/%s] to TiDB cluster [%s/%s]",
		from.Namespace, from.DrainerName, to.Namespace, to.ClusterName)

	// TODO: better incremental files restore solution
	reparoConfig := strings.Join([]string{
		`data-dir = "/data/pb"`,
		`log-level = "info"`,
		`dest-type = "mysql"`,
		`[dest-db]`,
		fmt.Sprintf(`host = "%s"`, util.GetTidbServiceName(to.ClusterName)),
		"port = 4000",
		`user = "root"`,
		`password = ""`,
	}, "\n")

	temp, err := template.New("reparo-command").Parse(RunReparoCommandTemplate)
	if err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	if err := temp.Execute(buff, &struct {
		ReparoConfig string
		PodName      string
	}{
		ReparoConfig: reparoConfig,
		PodName:      fmt.Sprintf("%s-%s-drainer-0", from.SourceClusterName, from.DrainerName),
	}); err != nil {
		return err
	}

	if res, err := exec.Command("/bin/sh", "-c", buff.String()).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restore incremental files from dainer [%s/%s] to TiDB cluster [%s/%s], %v, %s",
			from.Namespace, from.DrainerName, to.Namespace, to.ClusterName, err, res)
	}
	return nil
}

func (oa *operatorActions) CheckRestoreIncrementalFiles(from *DrainerConfig, to *TidbClusterConfig) error {
	return nil
}
