package message

import (
	"net/http"
	"time"

	"demo-server/lib/errors"
)

//  号段　　1000000 - 1999999
//  公共错误 1110000
//  业务错误 12XXXXX
//  为防止定义重复，以及方便查找，各个服务错误码可直接选用一个段号
//  调用服务出错已选用一个段号，参见最后
const (
	CodeOK = 200

	// 通用错误码

	// 接口错误
	CodeUnknowErr   = 1110000
	CodeParamErr    = 1110001
	CodeServerErr   = 1110002
	CodeTokenErr    = 1110003
	CodeOvsTokenErr = 1110004
	CodeSignErr     = 1110005

	// redis错误
	CodeRedisKeyExistErr = 1120000

	// mysql错误
	CodeMysqlScanErr            = 1130000
	CodeMysqlSelectErr          = 1130001
	CodeMysqlUpdateErr          = 1130002
	CodeMysqlInsertErr          = 1130003
	CodeMysqlFieldTagOrmErr     = 1130010
	CodeMysqlInsertManyLimitErr = 1130011

	// 业务错误码

	// 公共模块
	CodeCommonUnknowErr = 1210000
	// 安防错误
	CodeSecurityUnknowErr = 1220000
	// 智能家居错误
	CodeHomeUnknowErr = 1230000
	// 消息中心错误
	CodeNoticeUnknowErr  = 1240000
	CodeInvalidDeviceErr = 1240001
	// 用户中心错误
	CodeUserUnknowErr              = 1250000
	CodeUserLoginFailErr           = 1250001
	CodeUserSendCodeFailErr        = 1250002
	CodeUserPermissionCheckFailErr = 1250003
	CodeUserOffLineErr             = 1250004
	CodeUserInfoErr                = 1250005
	CodeRefreshOVTokenErr          = 1250006
	CodeRefreshLoginTokenErr       = 1250007
	// 机器人控制错误
	CodeUserLogOutFailErr  = 1250008
	CodeGetUserInfoFailErr = 1250009
	CodeGetOVTokenFailErr  = 1250010
	CodeBindCodeInvalidErr = 1250015
	// 家庭信息错误
	CodeFamilyInfoErr    = 1250011
	CodeMembersInfoErr   = 1250012
	CodeMemberDeleteErr  = 1250013
	CodeAdminTransferErr = 1250014

	//视觉接口错误
	CodeFaceSearchErr = 1250021

	CodeRobotUnknowErr     = 1260000
	CodeVideoUnknowErr     = 1270000
	CodeVideoMsgUnknowErr  = 1270001
	CodeVideoLinkUnknowErr = 1270002
	CodeVideoEngagedErr    = 1270003
	// 技能服务错误
	CodeSkillsUnknowErr = 1280000

	// 资源服务
	CodeResUnknownErr    = 1290000
	CodeResUserInfoErr   = 1290001
	CodeResRobotInfoErr  = 1290002
	CodeResTaskListErr   = 1290003
	CodeResMemberListErr = 1290004
	CodeResScheduleErr   = 1290005

	// IFTTT
	CodeIftttUnknownErr         = 1300000
	CodeIftttHisStatusAddErr    = 1300001
	CodeIftttHisStatusUpdateErr = 1300002

	// 地图
	CodeMapRobotUnknowErr = 1400404
	CodeMapPkgMD5Err      = 1400410
	CodeMapPkgSaveErr     = 1400411
	CodeMapPkgLoadErr     = 1400412
	CodeMapDetailErr      = 1400413
	CodeMapDeviceErr      = 1400414
	CodeMapSaveErr        = 1400500
	CodeMapSitesSaveErr   = 1400501
	CodeMapDeleteErr      = 1400502
	CodeMapListErr        = 1400503
	CodeMapSwitchErr      = 1400504
	CodeMapSyncErr        = 1400505

	// 状态上报
	CodeRobotStatusUnknowErr = 1500500
)

var codeHttpStatus = map[int]int{
	CodeOK: http.StatusOK,

	CodeMapPkgSaveErr:        http.StatusInternalServerError,
	CodeMapSaveErr:           http.StatusInternalServerError,
	CodeMapSitesSaveErr:      http.StatusInternalServerError,
	CodeMapDeleteErr:         http.StatusInternalServerError,
	CodeMapListErr:           http.StatusInternalServerError,
	CodeMapSwitchErr:         http.StatusInternalServerError,
	CodeRobotStatusUnknowErr: http.StatusInternalServerError,
}

var codeText = map[int]string{
	CodeOK:                         "OK",
	CodeParamErr:                   "param error",
	CodeServerErr:                  "server internal error",
	CodeTokenErr:                   "token invalid",
	CodeSignErr:                    "sign invalid",
	CodeRedisKeyExistErr:           "redis key is exist",
	CodeMysqlScanErr:               "mysql scan error",
	CodeMysqlSelectErr:             "mysql select error",
	CodeMysqlFieldTagOrmErr:        "add 'orm' tag to the struct",
	CodeMysqlInsertManyLimitErr:    "batch insertion limit 0-100",
	CodeUserUnknowErr:              "用户中心未知错误",
	CodeUserLoginFailErr:           "用户登陆失败",
	CodeUserSendCodeFailErr:        "验证码发送失败",
	CodeUserPermissionCheckFailErr: "权限校验失败",
	CodeUserOffLineErr:             "已离线",
	CodeUserLogOutFailErr:          "退出登录失败",
	CodeGetUserInfoFailErr:         "获取用户信息失败",
	CodeGetOVTokenFailErr:          "获取OVtoken失败",
	CodeBindCodeInvalidErr:         "绑定码失效",
	CodeInvalidDeviceErr:           "invalid Device-ID header",
	CodeUserInfoErr:                "获取用户信息失败",
	CodeRefreshOVTokenErr:          "刷新OV Token失败",
	CodeRefreshLoginTokenErr:       "刷新登录 Token失败",
	CodeOvsTokenErr:                "ovs token invalid",
	CodeFaceSearchErr:              "视觉查询错误",
	CodeHomeUnknowErr:              "iot接口错误",
	CodeRobotUnknowErr:             "robot control 接口错误",
	CodeVideoUnknowErr:             "vedio 接口错误",
	CodeVideoEngagedErr:            "对方忙线",
	CodeVideoMsgUnknowErr:          "控制指令发送失败",
	CodeVideoLinkUnknowErr:         "获取视频通道失败",
	CodeResUnknownErr:              "资源下发未知错误",
	CodeResUserInfoErr:             "获取用户信息失败",
	CodeResRobotInfoErr:            "获取机器人信息失败",
	CodeResMemberListErr:           "获取家庭成员列表失败",
	CodeResTaskListErr:             "获取任务列表失败",
	CodeResScheduleErr:             "获取日程列表失败",
	CodeFamilyInfoErr:              "获取家庭信息失败",
	CodeMembersInfoErr:             "获取家庭成员失败",
	CodeMemberDeleteErr:            "删除家庭成员失败",
	CodeAdminTransferErr:           "管理员权限转让失败",
	CodeIftttUnknownErr:            "IFTTT未知错误",
	CodeIftttHisStatusAddErr:       "IFTTT开始执行的任务新增失败",
	CodeIftttHisStatusUpdateErr:    "IFTTT任务状态更新失败",
	CodeMapRobotUnknowErr:          "地图模块找到不机器人",
	CodeMapPkgMD5Err:               "地图文件MD5不一致",
	CodeMapPkgSaveErr:              "地图文件保存失败",
	CodeMapPkgLoadErr:              "地图文件加载失败",
	CodeMapSaveErr:                 "地图保存失败",
	CodeMapSitesSaveErr:            "地图地点保存失败",
	CodeMapDeleteErr:               "地图删除失败",
	CodeMapListErr:                 "获取地图列表失败",
	CodeMapSwitchErr:               "地图切换失败",
	CodeMapDetailErr:               "获取地图详情失败",
	CodeMapDeviceErr:               "设备未激活",
	CodeMapSyncErr:                 "同步地图错误",
	CodeRobotStatusUnknowErr:       "状态上报错误",
}

func HttpStatus(code int) int {
	if status, ok := codeHttpStatus[code]; ok {
		return status
	}
	return http.StatusBadRequest
}

func Text(code int) string {
	return codeText[code]
}

type ResponseHeader struct {
	Msg       string `json:"msg"`
	Code      int    `json:"code"`
	Time      int64  `json:"time"`
	RequestID string `json:"request_id"`
}

type ResponseData struct {
	Header ResponseHeader `json:"header"`
	Data   interface{}    `json:"data"`
}

type responseNoData struct {
	Header struct {
		Msg       string `json:"msg"`
		Code      int    `json:"code"`
		Time      int64  `json:"time"`
		RequestID string `json:"request_id"`
	} `json:"header"`
}

func Response(code int, data interface{}, rid string) *ResponseData {
	return &ResponseData{
		Header: struct {
			Msg       string `json:"msg"`
			Code      int    `json:"code"`
			Time      int64  `json:"time"`
			RequestID string `json:"request_id"`
		}{Msg: Text(code), Code: code, Time: time.Now().Unix(), RequestID: rid},
		Data: data,
	}
}

func ResponseNoData(code int, rid string) *responseNoData {
	return &responseNoData{
		Header: struct {
			Msg       string `json:"msg"`
			Code      int    `json:"code"`
			Time      int64  `json:"time"`
			RequestID string `json:"request_id"`
		}{Msg: Text(code), Code: code, Time: time.Now().Unix(), RequestID: rid},
	}
}

type responseNoHeader struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func ResponseNoHeader(code int, data interface{}) *responseNoHeader {
	return &responseNoHeader{
		Msg:  Text(code),
		Code: code,
	}
}

func ServerErrorResponse(err error, rid string) *ResponseData {
	if err, ok := err.(*errors.Error); ok {
		return Response(err.Code, nil, rid)
	}
	return &ResponseData{
		Header: struct {
			Msg       string `json:"msg"`
			Code      int    `json:"code"`
			Time      int64  `json:"time"`
			RequestID string `json:"request_id"`
		}{Msg: err.Error(), Code: http.StatusInternalServerError, Time: time.Now().Unix(), RequestID: rid},
	}
}

type paginatedData struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

func PaginatedData(list interface{}, total int) *paginatedData {
	return &paginatedData{
		Total: total,
		List:  list,
	}
}
