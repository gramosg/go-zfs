package ioctl

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

func TestSequence(t *testing.T) {
	baseLocation := "/dev/shm"

	Init("")

	fileLocation := filepath.Join(baseLocation, "test.img")
	f, err := os.Create(fileLocation)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Truncate(1e9); err != nil { // 1GiB
		t.Fatal(err)
	}
	f.Close()

	defer PoolDestroy("tp1")
	defer os.Remove(fileLocation)

	err = PoolCreate("tp1", map[string]uint64{}, VDev{
		Type: "root",
		Children: []VDev{
			VDev{
				Type: "file",
				Path: fileLocation,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	poolConfigs, err := PoolConfigs()
	assert.NoError(t, err)
	assert.Contains(t, poolConfigs, "tp1")

	_, _, _, _, err = DatasetListNext("tp1", 0)
	if err != unix.ESRCH {
		t.Errorf("Dataset list of empty pool doesn't return ESRCH (instead %v)", err)
	}

	if err := Create("tp1/test5", ObjectTypeZFS, &DatasetProps{"mountpoint": "legacy"}); err != nil {
		t.Fatal(err)
	}
	if err := Create("tp1/test7", ObjectTypeZFS, &DatasetProps{"mountpoint": "legacy"}); err != nil {
		t.Error(err)
	}
	err = Create("tp1/test7/ml0", ObjectTypeZFS, &DatasetProps{"mountpoint": "legacy"})
	assert.NoError(t, err)

	props, err := ObjsetStats("tp1/test5")
	assert.NoError(t, err, "Failed to call ObjsetStats")

	name, cookie, _, props, err := DatasetListNext("tp1", 0)
	assert.NoError(t, err)
	assert.Equal(t, props["mountpoint"].Value.(string), "legacy")
	assert.Equal(t, props["type"].Value.(uint64), uint64(2))

	name2, cookie, _, props, err := DatasetListNext("tp1", cookie)
	assert.NoError(t, err)
	assert.NotEqual(t, name, name2) // Test if cookies work

	var findDatasetsRecursive func(prefix string, root bool) bool
	findDatasetsRecursive = func(prefix string, root bool) bool {
		cookie := uint64(0)
		for {
			var name string
			name, cookie, _, _, err = DatasetListNext(prefix, cookie)
			if err == unix.ESRCH {
				if root {
					t.Error("Didn't find ml0 in listing")
				}
				return false
			}
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			if strings.Contains(name, "ml0") {
				return true
			}
			if findDatasetsRecursive(name, false) {
				return true
			}
		}
	}

	findDatasetsRecursive("tp1", true)

	if err := Rename("tp1/test7", "tp1/test6", false); err != nil {
		t.Error(err)
	}
	if err := Snapshot([]string{"tp1/test5@snap1", "tp1/test6@snap1"}, "tp1", nil); err != nil {
		t.Error(err)
	}
	n, err := SendSpace("tp1/test5@snap1", SendSpaceOptions{Compress: true})
	if err != nil {
		t.Error(err)
	}
	if n == 0 {
		t.Error(errors.New("size of snaphsot is 0"))
	}
	if err := Clone("tp1/test5@snap1", "tp1/test9", nil); err != nil {
		t.Error(err)
	}
	if err := Snapshot([]string{"tp1/test5@snap2"}, "tp1", nil); err != nil {
		t.Error(err)
	}

	written, err := GetSpaceWritten("tp1/test5", "tp1/test5@snap2")
	assert.NoError(t, err, "GetSpaceWritten failed")
	assert.Zero(t, written, "written is not zero for a fresh snapshot")

	n, err = SendSpace("tp1/test5@snap2", SendSpaceOptions{From: "tp1/test5@snap1"})
	if err != nil {
		t.Error(err)
	}
	if n == 0 {
		t.Error(errors.New("size of snaphsot is 0"))
	}

	r, err := Send("tp1/test5@snap2", SendOptions{From: "tp1/test5@snap1"})
	assert.NoError(t, err, "Failed to send snap2")
	defer r.Close()

	sendLocation := filepath.Join(baseLocation, "send.bin")
	f, err = os.Create(sendLocation)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	defer os.Remove(sendLocation)
	if _, err := io.Copy(f, r); err != nil {
		t.Error(err)
	}

	r.Close()
	f.Seek(0, io.SeekStart)

	Destroy("tp1/test5@snap2", ObjectTypeAny, false)

	w, err := Receive("tp1", ReceiveOpts{SnapshotName: "tp1/test5@snap2"})
	assert.NoError(t, err, "Failed to start receive")
	_, err = io.Copy(w, f)
	assert.NoError(t, err, "Failed to receive data")
	err = w.WaitAndClose()
	assert.NoError(t, err, "Failed to complete receive")

	props, err = ObjsetStats("tp1/test5@snap2")
	assert.NoError(t, err, "Failed to stat received snapshot")
	assert.Equal(t, ObjectType(props["type"].Value.(uint64)), ObjectTypeZFS, "Invalid received object type")

	r, err = Send("tp1/test5@nonexistent", SendOptions{})
	if err == nil {
		t.Error("Nonexistent send should immediately return an error")
	}

	if err := StartStopScan("tp1", ScanTypeScrub); err != nil {
		t.Error(err)
	}

	// TODO: Look if scrub is running

	if err := RegenerateGUID("tp1"); err != nil {
		t.Error(err)
	}

	configs, err := PoolConfigs()
	if err != nil {
		t.Fatal(err)
	}

	if err := PoolExport("tp1", false, false); err != nil {
		t.Fatal(err)
	}
	importConfig := configs["tp1"].(map[string]interface{})

	if _, err := PoolImport("tp1", importConfig, nil); err != nil {
		t.Fatal(err)
	}

	// TODO: Validate that GUID has changed

	if err := Destroy("tp1/test9", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}
	if err := Destroy("tp1/test5@snap1", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}
	if err := Destroy("tp1/test5@snap2", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}
	if err := Destroy("tp1/test6@snap1", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}
	if err := Destroy("tp1/test6/ml0", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}

	if err := Destroy("tp1/test6", ObjectTypeAny, false); err != nil {
		t.Error(err)
	}
	if err := PoolDestroy("tp1"); err != nil {
		t.Error(err)
	}
}
