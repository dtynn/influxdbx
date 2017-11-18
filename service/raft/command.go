package raft

type fsmCmdType byte

const (
	fsmCmdTypeUnknown fsmCmdType = 0

	fsmCmdTypeCreateContinuousQuery = 11
	fsmCmdTypeDropContinuousQuery   = 12

	fsmCmdTypeCreateDatabase                    = 21
	fsmCmdTypeCreateDatabaseWithRetentionPolicy = 22
	fsmCmdTypeDropDatabase                      = 23

	fsmCmdTypeCreateRetentionPolicy = 31
	fsmCmdTypeDropRetentionPolicy   = 32
	fsmCmdTypeUpdateRetentionPolicy = 33

	fsmCmdTypeCreateSubscription = 41
	fsmCmdTypeDropSubscription   = 42

	fsmCmdTypeCreateUser        = 51
	fsmCmdTypeDropUser          = 52
	fsmCmdTypeSetAdminPrivilege = 53
	fsmCmdTypeSetPrivilege      = 54
	fsmCmdTypeUpdateUser        = 55

	fsmCmdTypeDropShard              = 61
	fsmCmdTypeShardGroupsByTimeRange = 62
)

type fsmCmdResponse struct {
	res interface{}
	err error
}

func newFsmCmdResponse(res interface{}, err error) fsmCmdResponse {
	return fsmCmdResponse{
		res: res,
		err: err,
	}
}
