package ioctl

type Cmd struct {
	Name              [4096]byte
	Nvlist_src        uint64
	Nvlist_src_size   uint64
	Nvlist_dst        uint64
	Nvlist_dst_size   uint64
	Nvlist_dst_filled bool
	Pad_cgo_0         [3]byte
	Pad2              int32
	History           uint64
	Value             [8192]int8
	String            [256]int8
	Guid              uint64
	Nvlist_conf       uint64
	Nvlist_conf_size  uint64
	Cookie            uint64
	Objset_type       uint64
	Perm_action       uint64
	History_len       uint64
	History_offset    uint64
	Obj               uint64
	Iflags            uint64
	Share             Share
	Objset_stats      DMUObjectSetStats
	Begin_record      DRRBegin
	Inject_record     ZInjectRecord
	Defer_destroy     uint32
	Flags             uint32
	Action_handle     uint64
	Cleanup_fd        int32
	Simple            uint8
	Pad               [3]uint8
	Sendobj           uint64
	Fromobj           uint64
	Createtxg         uint64
	Stat              Stat
}
type DMUObjectType uint32
type DMUObjectSetStats struct {
	Num_clones   uint64
	Creation_txg uint64
	Guid         uint64
	Type         uint32
	Is_snapshot  uint8
	Inconsistent uint8
	Origin       [256]int8
	Pad_cgo_0    [2]byte
}
type ZInjectRecord struct {
	Objset    uint64
	Object    uint64
	Start     uint64
	End       uint64
	Guid      uint64
	Level     uint32
	Error     uint32
	Type      uint64
	Freq      uint32
	Failfast  uint32
	Func      [255]int8
	Pad_cgo_0 [1]byte
	Iotype    uint32
	Duration  int32
	Timer     uint64
	Nlanes    uint64
	Cmd       uint32
	Pad       uint32
}
type Share struct {
	Exportdata uint64
	Sharedata  uint64
	Sharetype  uint64
	Sharemax   uint64
}
type Stat struct {
	Gen   uint64
	Mode  uint64
	Links uint64
	Ctime [2]uint64
}
type DRRBegin struct {
	Magic         uint64
	Versioninfo   uint64
	Creation_time uint64
	Type          uint32
	Flags         uint32
	Toguid        uint64
	Fromguid      uint64
	Toname        [255]int8
	Pad_cgo_0     [1]byte
}
type Ioctl uint32

const (
	ZFS_IOC_POOL_CREATE Ioctl = ('Z' << 8) + iota
	ZFS_IOC_POOL_DESTROY
	ZFS_IOC_POOL_IMPORT
	ZFS_IOC_POOL_EXPORT
	ZFS_IOC_POOL_CONFIGS
	ZFS_IOC_POOL_STATS
	ZFS_IOC_POOL_TRYIMPORT
	ZFS_IOC_POOL_SCAN
	ZFS_IOC_POOL_FREEZE
	ZFS_IOC_POOL_UPGRADE
	ZFS_IOC_POOL_GET_HISTORY
	ZFS_IOC_VDEV_ADD
	ZFS_IOC_VDEV_REMOVE
	ZFS_IOC_VDEV_SET_STATE
	ZFS_IOC_VDEV_ATTACH
	ZFS_IOC_VDEV_DETACH
	ZFS_IOC_VDEV_SETPATH
	ZFS_IOC_VDEV_SETFRU
	ZFS_IOC_OBJSET_STATS
	ZFS_IOC_OBJSET_ZPLPROPS
	ZFS_IOC_DATASET_LIST_NEXT
	ZFS_IOC_SNAPSHOT_LIST_NEXT
	ZFS_IOC_SET_PROP
	ZFS_IOC_CREATE
	ZFS_IOC_DESTROY
	ZFS_IOC_ROLLBACK
	ZFS_IOC_RENAME
	ZFS_IOC_RECV
	ZFS_IOC_SEND
	ZFS_IOC_INJECT_FAULT
	ZFS_IOC_CLEAR_FAULT
	ZFS_IOC_INJECT_LIST_NEXT
	ZFS_IOC_ERROR_LOG
	ZFS_IOC_CLEAR
	ZFS_IOC_PROMOTE
	ZFS_IOC_SNAPSHOT
	ZFS_IOC_DSOBJ_TO_DSNAME
	ZFS_IOC_OBJ_TO_PATH
	ZFS_IOC_POOL_SET_PROPS
	ZFS_IOC_POOL_GET_PROPS
	ZFS_IOC_SET_FSACL
	ZFS_IOC_GET_FSACL
	ZFS_IOC_SHARE
	ZFS_IOC_INHERIT_PROP
	ZFS_IOC_SMB_ACL
	ZFS_IOC_USERSPACE_ONE
	ZFS_IOC_USERSPACE_MANY
	ZFS_IOC_USERSPACE_UPGRADE
	ZFS_IOC_HOLD
	ZFS_IOC_RELEASE
	ZFS_IOC_GET_HOLDS
	ZFS_IOC_OBJSET_RECVD_PROPS
	ZFS_IOC_VDEV_SPLIT
	ZFS_IOC_NEXT_OBJ
	ZFS_IOC_DIFF
	ZFS_IOC_TMP_SNAPSHOT
	ZFS_IOC_OBJ_TO_STATS
	ZFS_IOC_SPACE_WRITTEN
	ZFS_IOC_SPACE_SNAPS
	ZFS_IOC_DESTROY_SNAPS
	ZFS_IOC_POOL_REGUID
	ZFS_IOC_POOL_REOPEN
	ZFS_IOC_SEND_PROGRESS
	ZFS_IOC_LOG_HISTORY
	ZFS_IOC_SEND_NEW
	ZFS_IOC_SEND_SPACE
	ZFS_IOC_CLONE
	ZFS_IOC_BOOKMARK
	ZFS_IOC_GET_BOOKMARKS
	ZFS_IOC_DESTROY_BOOKMARKS
	ZFS_IOC_RECV_NEW
	ZFS_IOC_POOL_SYNC

	/*
	 * Linux - 3/64 numbers reserved.
	 */
	ZFS_IOC_LINUX Ioctl = ('Z' << 8) + 0x80 + iota
	ZFS_IOC_EVENTS_NEXT
	ZFS_IOC_EVENTS_CLEAR
	ZFS_IOC_EVENTS_SEEK

	/*
	 * FreeBSD - 1/64 numbers reserved.
	 */
	ZFS_IOC_FREEBSD Ioctl = ('Z' << 8) + 0xC0 + iota

	ZFS_IOC_LAST
)
