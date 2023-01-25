package main

import (
	"encoding/json"
	"fmt"

	snapshotv1api "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
)

func main() {
	volumeSnapshotByte := []byte(`{"apiVersion":"snapshot.storage.k8s.io/v1","kind":"VolumeSnapshot","metadata":{"creationTimestamp":"2022-08-01T03:26:30Z","generation":1,"managedFields":[{"apiVersion":"snapshot.storage.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{".":{},"f:source":{".":{},"f:persistentVolumeClaimName":{}},"f:volumeSnapshotClassName":{}}},"manager":"aksdev","operation":"Update","time":"2022-08-01T03:26:30Z"}],"name":"csi-volume-snapshot","namespace":"e2e-ns-cbjkfi0sa298vhb5i7e0","resourceVersion":"3169","uid":"8f856981-4ea6-4110-a85a-7425f04809fc"},"spec":{"source":{"persistentVolumeClaimName":"pvc-mdkhk"},"volumeSnapshotClassName":"csi-volume-snapshot-class"}}`)
	volSnap := &snapshotv1api.VolumeSnapshot{}
	err := json.Unmarshal(volumeSnapshotByte, volSnap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", volSnap.Name)
}
