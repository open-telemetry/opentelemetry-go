// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setDefaultContainerProviders() {
	setContainerProviders(
		getContainerIDFromCGroup,
	)
}

func setContainerProviders(
	idProvider containerIDProvider,
) {
	containerID = idProvider
}

func TestGetContainerIDFromReader(t *testing.T) {
	testCases := []struct {
		name                string
		reader              io.Reader
		expectedContainerID string
	}{
		{
			name: "multiple lines",
			reader: strings.NewReader(`//
1:cpuset:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23
1:cpuset:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d24
`),
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name: "no container id",
			reader: strings.NewReader(`//
1:name=systemd:/podruntime/docker
`),
			expectedContainerID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromReader(tc.reader, getContainerIDFromCgroupV1Line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}

func TestGetContainerIDFromCGroup(t *testing.T) {
	t.Cleanup(func() {
		osStat = defaultOSStat
		osOpen = defaultOSOpen
	})

	testCases := []struct {
		name                 string
		cgroupV1FileNotExist bool
		cgroupV2FileNotExist bool
		openFileError        error
		cgroupV1FileContent  string
		cgroupV2FileContent  string
		expectedContainerID  string
		expectedError        bool
	}{
		{
			name:                 "the cgroup file does not exist",
			cgroupV1FileNotExist: true,
			cgroupV2FileNotExist: true,
		},
		{
			name:          "error when opening cgroup file",
			openFileError: errors.New("test"),
			expectedError: true,
		},
		{
			name:                "cgroup file v1",
			cgroupV1FileContent: "1:cpuset:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23",
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "cgroup file v2",
			cgroupV2FileContent: "474 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/hostname /etc/hosts rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
		{
			name:                "both way fail",
			cgroupV1FileContent: " ",
			cgroupV2FileContent: " ",
		},
		{
			name: "minikube containerd",
			cgroupV1FileContent: `11:cpuset:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
10:hugetlb:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
9:pids:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
8:memory:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
7:net_cls,net_prio:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
6:perf_event:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
5:blkio:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
4:devices:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
3:freezer:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
2:cpu,cpuacct:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
1:name=systemd:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236`,
			cgroupV2FileContent: `1517 1428 0:208 / / rw,relatime master:510 - overlay overlay rw,lowerdir=/mnt/sda1/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/34/fs,upperdir=/mnt/sda1/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/35/fs,workdir=/mnt/sda1/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/35/work
1518 1517 0:210 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
1519 1517 0:211 / /dev rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
1520 1519 0:212 / /dev/pts rw,nosuid,noexec,relatime - devpts devpts rw,gid=5,mode=620,ptmxmode=666
1521 1519 0:198 / /dev/mqueue rw,nosuid,nodev,noexec,relatime - mqueue mqueue rw
1522 1517 0:203 / /sys ro,nosuid,nodev,noexec,relatime - sysfs sysfs ro
1523 1522 0:213 / /sys/fs/cgroup rw,nosuid,nodev,noexec,relatime - tmpfs tmpfs rw,mode=755
1524 1523 0:24 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:8 - cgroup cgroup rw,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd
1525 1523 0:26 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/cpu,cpuacct ro,nosuid,nodev,noexec,relatime master:11 - cgroup cgroup rw,cpu,cpuacct
1526 1523 0:27 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/freezer ro,nosuid,nodev,noexec,relatime master:12 - cgroup cgroup rw,freezer
1527 1523 0:28 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/devices ro,nosuid,nodev,noexec,relatime master:13 - cgroup cgroup rw,devices
1528 1523 0:29 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/blkio ro,nosuid,nodev,noexec,relatime master:14 - cgroup cgroup rw,blkio
1529 1523 0:30 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/perf_event ro,nosuid,nodev,noexec,relatime master:15 - cgroup cgroup rw,perf_event
1530 1523 0:31 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/net_cls,net_prio ro,nosuid,nodev,noexec,relatime master:16 - cgroup cgroup rw,net_cls,net_prio
1531 1523 0:32 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/memory ro,nosuid,nodev,noexec,relatime master:17 - cgroup cgroup rw,memory
1532 1523 0:33 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/pids ro,nosuid,nodev,noexec,relatime master:18 - cgroup cgroup rw,pids
1533 1523 0:34 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/hugetlb ro,nosuid,nodev,noexec,relatime master:19 - cgroup cgroup rw,hugetlb
1534 1523 0:35 /kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236 /sys/fs/cgroup/cpuset ro,nosuid,nodev,noexec,relatime master:20 - cgroup cgroup rw,cpuset
1535 1517 8:1 /var/lib/kubelet/pods/28478e30-384f-41e5-9d85-eae249ae8506/etc-hosts /etc/hosts rw,relatime - ext4 /dev/sda1 rw
1536 1519 8:1 /var/lib/kubelet/pods/28478e30-384f-41e5-9d85-eae249ae8506/containers/alpine/e3d5dec7 /dev/termination-log rw,relatime - ext4 /dev/sda1 rw
1537 1517 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw
1538 1517 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6/resolv.conf /etc/resolv.conf rw,relatime - ext4 /dev/sda1 rw
1539 1519 0:195 / /dev/shm rw,nosuid,nodev,noexec,relatime - tmpfs shm rw,size=65536k
1540 1517 0:194 / /run/secrets/kubernetes.io/serviceaccount ro,relatime - tmpfs tmpfs rw,size=5925720k
1429 1519 0:212 /0 /dev/console rw,nosuid,noexec,relatime - devpts devpts rw,gid=5,mode=620,ptmxmode=666
1430 1518 0:210 /asound /proc/asound ro,nosuid,nodev,noexec,relatime - proc proc rw
1431 1518 0:210 /bus /proc/bus ro,nosuid,nodev,noexec,relatime - proc proc rw
1432 1518 0:210 /fs /proc/fs ro,nosuid,nodev,noexec,relatime - proc proc rw
1433 1518 0:210 /irq /proc/irq ro,nosuid,nodev,noexec,relatime - proc proc rw
1434 1518 0:210 /sys /proc/sys ro,nosuid,nodev,noexec,relatime - proc proc rw
1435 1518 0:210 /sysrq-trigger /proc/sysrq-trigger ro,nosuid,nodev,noexec,relatime - proc proc rw
1436 1518 0:214 / /proc/acpi ro,relatime - tmpfs tmpfs ro
1437 1518 0:211 /null /proc/kcore rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
1438 1518 0:211 /null /proc/keys rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
1439 1518 0:211 /null /proc/timer_list rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
1440 1518 0:215 / /proc/scsi ro,relatime - tmpfs tmpfs ro
1441 1522 0:216 / /sys/firmware ro,relatime - tmpfs tmpfs ro`,
			expectedContainerID: "58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236",
		},
		{
			name: "minikube docker cgroup",
			cgroupV1FileContent: `11:blkio:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
10:perf_event:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
9:pids:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
8:memory:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
7:cpu,cpuacct:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
6:net_cls,net_prio:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
5:cpuset:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
4:devices:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
3:hugetlb:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
2:freezer:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
1:name=systemd:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope`,
			expectedContainerID: "3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b",
		},
		{
			name: "podman cgroup",
			cgroupV1FileContent: `14:name=systemd:/user.slice/user-1000.slice/user@1000.service/app.slice/podman.service
13:rdma:/
12:pids:/user.slice/user-1000.slice/user@1000.service
11:hugetlb:/
10:net_prio:/
9:perf_event:/
8:net_cls:/
7:freezer:/
6:devices:/user.slice
5:blkio:/user.slice
4:cpuacct:/
3:cpu:/user.slice
2:cpuset:/
1:memory:/user.slice/user-1000.slice/user@1000.service
0::/user.slice/user-1000.slice/user@1000.service/app.slice/podman.service`,
			cgroupV2FileContent: `1088 875 0:118 / / rw,noatime - fuse.fuse-overlayfs fuse-overlayfs rw,user_id=0,group_id=0,default_permissions,allow_other
1089 1088 0:121 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
1090 1088 0:122 / /dev rw,nosuid,noexec - tmpfs tmpfs rw,size=65536k,mode=755,uid=1000,gid=1000
1091 1088 0:123 / /sys ro,nosuid,nodev,noexec,relatime - sysfs sysfs rw
1092 1090 0:124 / /dev/pts rw,nosuid,noexec,relatime - devpts devpts rw,gid=100004,mode=620,ptmxmode=666
1093 1090 0:120 / /dev/mqueue rw,nosuid,nodev,noexec,relatime - mqueue mqueue rw
1094 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/hosts /etc/hosts rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000
1095 1090 0:117 / /dev/shm rw,nosuid,nodev,noexec,relatime - tmpfs shm rw,size=64000k,uid=1000,gid=1000
1096 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/hostname /etc/hostname rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000
1097 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/.containerenv /run/.containerenv rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000
1098 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/run/secrets /run/secrets rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000
1099 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/resolv.conf /etc/resolv.conf rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000
1100 1091 0:125 / /sys/fs/cgroup rw,nosuid,nodev,noexec,relatime - tmpfs cgroup rw,size=1024k,uid=1000,gid=1000`,
			expectedContainerID: "1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = func(name string) (os.FileInfo, error) {
				if tc.cgroupV1FileNotExist && name == cgroupV1Path {
					return nil, os.ErrNotExist
				}
				if tc.cgroupV2FileNotExist && name == cgroupV2Path {
					return nil, os.ErrNotExist
				}
				return nil, nil
			}

			osOpen = func(name string) (io.ReadCloser, error) {
				if tc.openFileError != nil {
					return nil, tc.openFileError
				}
				if name == cgroupV1Path {
					return io.NopCloser(strings.NewReader(tc.cgroupV1FileContent)), nil
				}
				return io.NopCloser(strings.NewReader(tc.cgroupV2FileContent)), nil
			}

			containerID, err := getContainerIDFromCGroup()
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
