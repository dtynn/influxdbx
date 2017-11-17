package raft

type fsmCmdType byte

const (
	fsmCmdTypeUnknown fsmCmdType = 0

	fsmCmdTypeCreateContinuousQuery = 11
	fsmCmdTypeDropContinuousQuery   = 12

	fsmCmdTypeCreateDatabase                    = 21
	fsmCmdTypeCreateDatabaseWithRetentionPolicy = 22
	fsmCmdTypeDatabase                          = 23
	fsmCmdTypeDatabases                         = 24
	fsmCmdTypeDropDatabase                      = 25

	fsmCmdTypeCreateRetentionPolicy = 31
	fsmCmdTypeDropRetentionPolicy   = 32
	fsmCmdTypeRetentionPolicy       = 33
	fsmCmdTypeUpdateRetentionPolicy = 34

	fsmCmdTypeCreateSubscription = 41
	fsmCmdTypeDropSubscription   = 42

	fsmCmdTypeCreateUser        = 51
	fsmCmdTypeDropUser          = 52
	fsmCmdTypeSetAdminPrivilege = 53
	fsmCmdTypeSetPrivilege      = 54
	fsmCmdTypeUpdateUser        = 55
	fsmCmdTypeUserPrivilege     = 56
	fsmCmdTypeUserPrivileges    = 57
	fsmCmdTypeUsers             = 58

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
