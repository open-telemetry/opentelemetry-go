// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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

func TestGetContainerIDFromLine(t *testing.T) {
	testCases := []struct {
		name                string
		line                string
		expectedContainerID string
	}{
		{
			name:                "with suffix",
			line:                "13:name=systemd:/podruntime/docker/kubepods/ac679f8a8319c8cf7d38e1adf263bc08d23.aaaa",
			expectedContainerID: "ac679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "with prefix and suffix",
			line:                "13:name=systemd:/podruntime/docker/kubepods/crio-dc679f8a8319c8cf7d38e1adf263bc08d23.stuff",
			expectedContainerID: "dc679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "no prefix and suffix",
			line:                "13:name=systemd:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
			expectedContainerID: "d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
		},
		{
			name:                "with space",
			line:                " 13:name=systemd:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356 ",
			expectedContainerID: "d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
		},
		{
			name:                "with colon",
			line:                " 13:name=systemd:/kuberuntime/containerd/kubepods-pod872d2066_00ef_48ea_a7d8_51b18b72d739:cri-containerd:e857a4bf05a69080a759574949d7a0e69572e27647800fa7faff6a05a8332aa1",
			expectedContainerID: "e857a4bf05a69080a759574949d7a0e69572e27647800fa7faff6a05a8332aa1",
		},
		{
			name: "invalid hex string",
			line: "13:name=systemd:/podruntime/docker/kubepods/ac679f8a8319c8cf7d38e1adf263bc08d23zzzz",
		},
		{
			name: "no container id - 1",
			line: "pids: /",
		},
		{
			name: "no container id - 2",
			line: "pids: ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromLine(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
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
1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23
1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d24
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
			containerID := getContainerIDFromReader(tc.reader)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}

func TestGetContainerIDFromMountInfoLine(t *testing.T) {
	testCases := []struct {
		name                string
		line                string
		expectedContainerID string
	}{
		{
			name:                "crio prefix",
			line:                "7282 7281 0:27 /kubepods.slice/kubepods-burstable.slice/kubepods-burstable-pod8f215fa2_6177_4ab9_b1f4_c802d19657bc.slice/crio-f23ec1d4b715c6531a17e9c549222fbbe1f7ffff697a29a2212b3b4cdc37f52e.scope /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:9 - cgroup cgroup rw",
			expectedContainerID: "f23ec1d4b715c6531a17e9c549222fbbe1f7ffff697a29a2212b3b4cdc37f52e",
		},
		{
			name:                "cri-containerd prefix",
			line:                "2009 2008 0:32 /system.slice/containerd.service/kubepods-burstable-pod321c09bf_282b_44e4_a467_39daf144ef1f.slice:cri-containerd:f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91 /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:11 - cgroup cgroup rw,xattr,name=systemd",
			expectedContainerID: "f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91",
		},
		{
			name: "non-container line",
			line: "457 456 0:146 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw",
		},
		{
			name: "crio prefix with invalid hex",
			line: "100 99 0:27 /kubepods.slice/crio-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz.scope /sys/fs/cgroup/systemd",
		},
		{
			name: "crio prefix too short",
			line: "100 99 0:27 /kubepods.slice/crio-abc123.scope /sys/fs/cgroup/systemd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := getContainerIDFromMountInfoLine(tc.line)
			assert.Equal(t, tc.expectedContainerID, id)
		})
	}
}

func TestGetLastHexSegment(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "docker containers path",
			line:     "473 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expected: "be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2",
		},
		{
			name:     "docker in minikube",
			line:     "929 920 254:1 /docker/volumes/minikube/_data/lib/docker/containers/0eaa6718003210b6520f7e82d14b4c8d4743057a958a503626240f8d1900bc33/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expected: "0eaa6718003210b6520f7e82d14b4c8d4743057a958a503626240f8d1900bc33",
		},
		{
			name:     "podman overlay-containers",
			line:     "983 961 0:56 /containers/overlay-containers/2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6/userdata/hostname /etc/hostname ro - tmpfs tmpfs rw",
			expected: "2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6",
		},
		{
			name:     "containerd minikube sandboxes",
			line:     "1537 1517 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw",
			expected: "fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6",
		},
		{
			name: "no 64-hex segment",
			line: "100 99 0:50 /some/path/hostname /etc/hostname rw - ext4 /dev/sda1 rw",
		},
		{
			name: "invalid hex chars",
			line: "100 99 0:50 /containerd/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3fz9320f4402ae6/hostname /etc/hostname rw - ext4 /dev/sda1 rw",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := getLastHexSegment(tc.line)
			assert.Equal(t, tc.expected, id)
		})
	}
}

func TestGetContainerIDFromMountInfoReader(t *testing.T) {
	testCases := []struct {
		name                string
		content             string
		expectedContainerID string
	}{
		{
			name: "docker multi-line",
			content: `456 375 0:143 / / rw,relatime master:175 - overlay overlay rw,lowerdir=/var/lib/docker/overlay2/l/CBPR2ETR4Z3UMOOGIIRDVT2P27
457 456 0:146 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
466 456 0:147 / /dev rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
472 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/resolv.conf /etc/resolv.conf rw,relatime - ext4 /dev/vda1 rw
473 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw
474 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/hosts /etc/hosts rw,relatime - ext4 /dev/vda1 rw
377 457 0:146 /bus /proc/bus ro,nosuid,nodev,noexec,relatime - proc proc rw`,
			expectedContainerID: "be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2",
		},
		{
			name: "containerd with cri-containerd prefix returns workload ID not sandbox ID",
			content: `2002 1895 0:226 / / rw,relatime master:629 - overlay overlay rw
2009 2008 0:32 /system.slice/containerd.service/kubepods-burstable-pod321c09bf_282b_44e4_a467_39daf144ef1f.slice:cri-containerd:f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91 /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:11 - cgroup cgroup rw,xattr,name=systemd
2023 2002 253:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/b136f3d296b4c2024b3e7ad816f2a804a47cf1acc3d445075c6d78cf159ef58d/hostname /etc/hostname rw,relatime - xfs /dev/mapper/ubuntu--vg-root rw`,
			expectedContainerID: "f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91",
		},
		{
			name: "containerd prefix wins even when hostname line comes first",
			content: `2023 2002 253:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/b136f3d296b4c2024b3e7ad816f2a804a47cf1acc3d445075c6d78cf159ef58d/hostname /etc/hostname rw,relatime - xfs /dev/mapper/ubuntu--vg-root rw
2009 2008 0:32 /system.slice/containerd.service/kubepods-burstable-pod321c09bf_282b_44e4_a467_39daf144ef1f.slice:cri-containerd:f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91 /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:11 - cgroup cgroup rw,xattr,name=systemd`,
			expectedContainerID: "f2a44bc8e090f93a2b4d7f510bdaff0615ad52906e3287ee956dcf5aa5012a91",
		},
		{
			name: "crio with prefix returns workload ID",
			content: `7276 6904 0:507 / / rw,relatime - overlay overlay rw
7282 7281 0:27 /kubepods.slice/kubepods-burstable.slice/kubepods-burstable-pod8f215fa2_6177_4ab9_b1f4_c802d19657bc.slice/crio-f23ec1d4b715c6531a17e9c549222fbbe1f7ffff697a29a2212b3b4cdc37f52e.scope /sys/fs/cgroup/systemd ro,nosuid,nodev,noexec,relatime master:9 - cgroup cgroup rw
7304 7276 0:25 /containers/storage/overlay-containers/757a1c14bdd68b907c41f15436c0c2f9ec5a4cd4317135fcc1c4a64188db98d0/userdata/hostname /etc/hostname rw,nosuid,nodev master:28 - tmpfs tmpfs rw`,
			expectedContainerID: "f23ec1d4b715c6531a17e9c549222fbbe1f7ffff697a29a2212b3b4cdc37f52e",
		},
		{
			name: "crio generic without prefix uses last containers match over hostname",
			content: `10303 9025 0:676 / / rw,relatime master:2633 - overlay overlay rw
10312 10303 0:25 /containers/storage/overlay-containers/2ac4c84cb0d3c3beb04beeef6ccf71c17b5fdd0252ce3a2b66bc2fdd0aaa1814/userdata/hostname /etc/hostname rw,nosuid,nodev master:15 - tmpfs tmpfs rw
10314 10303 252:4 /ostree/deploy/rhcos/var/lib/kubelet/pods/0a947273-7214-4824-8411-875ebd7626e4/etc-hosts /etc/hosts rw,relatime - xfs /dev/vda4 rw
10316 10303 0:25 /containers/storage/overlay-containers/a8f62e52ed7c2cd85242dcf0eb1d727b643540ceca7f328ad7d2f31aedf07731/userdata/run/secrets /run/secrets rw,nosuid,nodev - tmpfs tmpfs rw`,
			expectedContainerID: "a8f62e52ed7c2cd85242dcf0eb1d727b643540ceca7f328ad7d2f31aedf07731",
		},
		{
			name: "containerd minikube with hostname only",
			content: `1239 872 0:60 / / rw,relatime master:451 - overlay overlay rw
1271 1239 0:62 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
1537 1517 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw
873 1271 0:62 /bus /proc/bus ro,nosuid,nodev,noexec,relatime - proc proc rw`,
			expectedContainerID: "fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6",
		},
		{
			name: "podman multi-line",
			content: `961 812 0:58 / / ro,relatime - overlay overlay rw,lowerdir=/home/dracula/.local/share/containers/storage/overlay/l/4NB35A5Z4YGWDHXYEUZU4FN6BU
962 961 0:63 / /sys ro,nosuid,nodev,noexec,relatime - sysfs sysfs rw
983 961 0:56 /containers/overlay-containers/2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6/userdata/hostname /etc/hostname ro,nosuid,nodev,noexec,relatime - tmpfs tmpfs rw,size=783888k`,
			expectedContainerID: "2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6",
		},
		{
			name: "gke containerd with hostname fallback only",
			content: `7634 6917 0:1147 / / rw,relatime - overlay overlay rw,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/673/fs:/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/671/fs:/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/669/fs:/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/668/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/915/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.gcfs/snapshotter/snapshots/915/work,uuid=on
7636 7634 0:1179 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
7637 7634 0:1180 / /dev rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
7638 7637 0:1181 / /dev/pts rw,nosuid,noexec,relatime - devpts devpts rw,gid=5,mode=620,ptmxmode=666
7639 7637 0:1149 / /dev/mqueue rw,nosuid,nodev,noexec,relatime - mqueue mqueue rw
7640 7634 0:1154 / /sys ro,nosuid,nodev,noexec,relatime - sysfs sysfs ro
7641 7640 0:27 / /sys/fs/cgroup ro,nosuid,nodev,noexec,relatime - cgroup2 cgroup rw
7642 7634 8:1 /var/lib/kubelet/pods/6bcf846c-5fe3-45cb-995f-5d34aa16e4c3/etc-hosts /etc/hosts rw,relatime - ext4 /dev/sda1 rw,commit=30
7643 7637 8:1 /var/lib/kubelet/pods/6bcf846c-5fe3-45cb-995f-5d34aa16e4c3/containers/latency-router/fc653337 /dev/termination-log rw,relatime - ext4 /dev/sda1 rw,commit=30
7644 7634 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/cbd4283245f21968f5a0f2960d4b689d063a0ccd27896b3665cfcbd140f3b5b5/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw,commit=30
7645 7634 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/cbd4283245f21968f5a0f2960d4b689d063a0ccd27896b3665cfcbd140f3b5b5/resolv.conf /etc/resolv.conf rw,relatime - ext4 /dev/sda1 rw,commit=30
7646 7637 0:947 / /dev/shm rw,relatime - tmpfs shm rw,size=65536k
7647 7634 0:702 / /run/secrets/kubernetes.io/serviceaccount ro,relatime - tmpfs tmpfs rw,size=28928128k,noswap
6918 7636 0:1179 /bus /proc/bus ro,nosuid,nodev,noexec,relatime - proc proc rw
6919 7636 0:1179 /fs /proc/fs ro,nosuid,nodev,noexec,relatime - proc proc rw
6920 7636 0:1179 /irq /proc/irq ro,nosuid,nodev,noexec,relatime - proc proc rw
6921 7636 0:1179 /sys /proc/sys ro,nosuid,nodev,noexec,relatime - proc proc rw
6922 7636 0:1179 /sysrq-trigger /proc/sysrq-trigger ro,nosuid,nodev,noexec,relatime - proc proc rw
6923 7636 0:1182 / /proc/acpi ro,relatime - tmpfs tmpfs ro
6924 7636 0:1180 /null /proc/interrupts rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
6925 7636 0:1180 /null /proc/kcore rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
6926 7636 0:1180 /null /proc/keys rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
6927 7636 0:1180 /null /proc/timer_list rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
6928 7636 0:1183 / /proc/scsi ro,relatime - tmpfs tmpfs ro
6929 7640 0:1184 / /sys/firmware ro,relatime - tmpfs tmpfs ro`,
			expectedContainerID: "cbd4283245f21968f5a0f2960d4b689d063a0ccd27896b3665cfcbd140f3b5b5",
		},
		{
			name: "orbstack docker multi-line",
			content: `291 219 0:76 / / rw,relatime - overlay overlay rw,lowerdir=/var/lib/docker/overlay2/l/EHOYRZ2XVF72EER2KUZMZRQ4AR:/var/lib/docker/overlay2/l/CPKJ4KNWRQ5UBZGO2NLTOWUR5H:/var/lib/docker/overlay2/l/LOSHS5MI4YOOHYHWFCMSYRMMOL:/var/lib/docker/overlay2/l/QGCTP7MCIDJDFW36CEJVHSCUBQ:/var/lib/docker/overlay2/l/NXOVYOW5ATJQZKPSKXAVF35YXR:/var/lib/docker/overlay2/l/7XE5L5WBX7ARM5KWBH7NBKOE5V:/var/lib/docker/overlay2/l/UVNTWA347GTLCKLFHLIF75TKX3:/var/lib/docker/overlay2/l/HCXXIKSKZIO7WAR4A3YZRMDERW,upperdir=/var/lib/docker/overlay2/e8ba3edcecf7d04bd9801d8085e49980f36a49a9b27b51b000ba4f423106543b/diff,workdir=/var/lib/docker/overlay2/e8ba3edcecf7d04bd9801d8085e49980f36a49a9b27b51b000ba4f423106543b/work
293 291 0:81 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
294 291 0:82 / /dev rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
295 294 0:83 / /dev/pts rw,nosuid,noexec,relatime - devpts devpts rw,gid=5,mode=620,ptmxmode=666
296 291 0:84 / /sys ro,nosuid,nodev,noexec,relatime - sysfs sysfs ro
297 296 0:31 / /sys/fs/cgroup ro,nosuid,nodev,noexec,relatime - cgroup2 cgroup rw,nsdelegate
298 294 0:79 / /dev/mqueue rw,nosuid,nodev,noexec,relatime - mqueue mqueue rw
299 294 0:85 / /dev/shm rw,nosuid,nodev,noexec,relatime - tmpfs shm rw,size=8207360k
300 291 0:37 /docker/volumes/test-depends_redis-data/_data /data rw,noatime master:65 - btrfs /dev/vdb1 rw,nodatasum,nodatacow,ssd,discard,space_cache=v2,subvolid=256,subvol=/docker
301 291 0:37 /docker/containers/6663b89b397dca3b1bc842cbeb31c5e3b1f27ec40ecb21ef15727300816e5158/resolv.conf /etc/resolv.conf rw,noatime - btrfs /dev/vdb1 rw,nodatasum,nodatacow,ssd,discard,space_cache=v2,subvolid=256,subvol=/docker
302 291 0:37 /docker/containers/6663b89b397dca3b1bc842cbeb31c5e3b1f27ec40ecb21ef15727300816e5158/hostname /etc/hostname rw,noatime - btrfs /dev/vdb1 rw,nodatasum,nodatacow,ssd,discard,space_cache=v2,subvolid=256,subvol=/docker
303 291 0:37 /docker/containers/6663b89b397dca3b1bc842cbeb31c5e3b1f27ec40ecb21ef15727300816e5158/hosts /etc/hosts rw,noatime - btrfs /dev/vdb1 rw,nodatasum,nodatacow,ssd,discard,space_cache=v2,subvolid=256,subvol=/docker
222 293 0:81 /bus /proc/bus ro,nosuid,nodev,noexec,relatime - proc proc rw
223 293 0:81 /fs /proc/fs ro,nosuid,nodev,noexec,relatime - proc proc rw
224 293 0:81 /irq /proc/irq ro,nosuid,nodev,noexec,relatime - proc proc rw
225 293 0:81 /sys /proc/sys ro,nosuid,nodev,noexec,relatime - proc proc rw
226 293 0:81 /sysrq-trigger /proc/sysrq-trigger ro,nosuid,nodev,noexec,relatime - proc proc rw
227 293 0:82 /null /proc/interrupts rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
228 293 0:82 /null /proc/keys rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
229 293 0:82 /null /proc/timer_list rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755
230 296 0:86 / /sys/firmware ro,relatime - tmpfs tmpfs ro`,
			expectedContainerID: "6663b89b397dca3b1bc842cbeb31c5e3b1f27ec40ecb21ef15727300816e5158",
		},
		{
			name: "no container id",
			content: `25 1 0:23 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
26 1 0:24 / /sys rw,nosuid,nodev,noexec,relatime - sysfs sysfs rw`,
			expectedContainerID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := getContainerIDFromMountInfoReader(strings.NewReader(tc.content))
			assert.Equal(t, tc.expectedContainerID, id)
		})
	}
}

func TestGetContainerIDFromMountInfo(t *testing.T) {
	t.Cleanup(func() {
		osStat = defaultOSStat
		osOpen = defaultOSOpen
	})

	testCases := []struct {
		name                string
		fileNotExist        bool
		openFileError       error
		content             string
		expectedContainerID string
		expectedError       bool
	}{
		{
			name:         "mountinfo file does not exist",
			fileNotExist: true,
		},
		{
			name:          "error when opening mountinfo file",
			openFileError: errors.New("test"),
			expectedError: true,
		},
		{
			name:                "mountinfo file with docker content",
			content:             "473 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2",
		},
		{
			name:    "mountinfo file with no container id",
			content: "25 1 0:23 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = func(string) (os.FileInfo, error) {
				if tc.fileNotExist {
					return nil, os.ErrNotExist
				}
				return nil, nil
			}

			osOpen = func(string) (io.ReadCloser, error) {
				if tc.openFileError != nil {
					return nil, tc.openFileError
				}
				return io.NopCloser(strings.NewReader(tc.content)), nil
			}

			containerID, err := getContainerIDFromMountInfo()
			assert.Equal(t, tc.expectedError, err != nil)
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
		cgroupFileNotExist   bool
		cgroupOpenError      error
		cgroupContent        string
		mountInfoFileNotExist bool
		mountInfoOpenError   error
		mountInfoContent     string
		expectedContainerID  string
		expectedError        bool
	}{
		{
			name:                  "neither file exists",
			cgroupFileNotExist:    true,
			mountInfoFileNotExist: true,
		},
		{
			name:            "error when opening cgroup file",
			cgroupOpenError: errors.New("test"),
			expectedError:   true,
		},
		{
			name:                "cgroup v1 has container id",
			cgroupContent:       "1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23",
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                 "cgroup v1 empty falls back to mountinfo",
			cgroupContent:        "0::/",
			mountInfoContent:     "473 456 254:1 /docker/containers/be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID:  "be522444b60caf2d3934b8b24b916a8a314f4b68d4595aa419874657e8d103f2",
		},
		{
			name:                 "cgroup file does not exist falls back to mountinfo",
			cgroupFileNotExist:   true,
			mountInfoContent:     "983 961 0:56 /containers/overlay-containers/2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6/userdata/hostname /etc/hostname ro - tmpfs tmpfs rw",
			expectedContainerID:  "2a33efc76e519c137fe6093179653788bed6162d4a15e5131c8e835c968afbe6",
		},
		{
			name:                  "cgroup v1 empty and mountinfo does not exist",
			cgroupContent:         "0::/",
			mountInfoFileNotExist: true,
		},
		{
			name:               "error when opening mountinfo file",
			cgroupContent:      "0::/",
			mountInfoOpenError: errors.New("test"),
			expectedError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = func(name string) (os.FileInfo, error) {
				switch name {
				case cgroupPath:
					if tc.cgroupFileNotExist {
						return nil, os.ErrNotExist
					}
				case mountInfoPath:
					if tc.mountInfoFileNotExist {
						return nil, os.ErrNotExist
					}
				}
				return nil, nil
			}

			osOpen = func(name string) (io.ReadCloser, error) {
				switch name {
				case cgroupPath:
					if tc.cgroupOpenError != nil {
						return nil, tc.cgroupOpenError
					}
					return io.NopCloser(strings.NewReader(tc.cgroupContent)), nil
				case mountInfoPath:
					if tc.mountInfoOpenError != nil {
						return nil, tc.mountInfoOpenError
					}
					return io.NopCloser(strings.NewReader(tc.mountInfoContent)), nil
				}
				return nil, os.ErrNotExist
			}

			containerID, err := getContainerIDFromCGroup()
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
