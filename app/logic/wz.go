package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"taogin/app/types"
	"taogin/app/types/pb"
	"taogin/config/global"
	"taogin/core/ws"
	"taogin/protobuf/protopb"
	"time"
)

type Wz struct {
}

func NewWz() *Wz {
	return &Wz{}
}

type Chess struct {
	I int32 `json:"i"`
	J int32 `json:"j"`
}

func (this *Wz) OnlineList(req *protopb.OnlineListReq, resp *protopb.OnlineListResp) error {
	var userList []*protopb.UserInfo
	for _, v := range ws.ClientManages.Clients {
		//过滤
		if len(req.Nickname) != 0 && req.Nickname == v.User.NickName {
			continue
		}
		fmt.Println(v.User.Sex)
		userList = append(userList, &protopb.UserInfo{
			Id:       v.UserId,
			Nickname: v.User.NickName,
			Avatar:   v.User.Avatar,
			//Sex:      v.User.Sex,
		})
	}
	resp.UserList = userList
	return nil
}

func (this *Wz) ApplyVs(req *protopb.ApplyVsReq, resp *protopb.ApplyVsResp, pbreq *pb.PbReq, cmd uint32, reqUid int64) error {
	if req.TargetUserId == pbreq.UserId {
		return errors.New("You can't play against yourself")
	}
	//验证状态
	wzUser, err := this.GetUserCacheById(req.TargetUserId)
	if err != nil {
		return errors.New(err.Error())
	}
	if wzUser.Status != 0 {
		return errors.New("VS failed")
	}

	//my user info
	userInfo := protopb.UserInfo{
		Id:       pbreq.User.Id,
		Nickname: pbreq.User.NickName,
		Avatar:   pbreq.User.Avatar,
	}
	//target user info
	//targetClient := ws.GetClient(req.TargetUserId)
	targetUserInfo := protopb.UserInfo{
		Id:       wzUser.User.Id,
		Nickname: wzUser.User.NickName,
		Avatar:   wzUser.User.Avatar,
	}
	//wzid
	room := userInfo.Id
	//my resp
	resp.UserInfo = &userInfo
	resp.TargetUserInfo = &targetUserInfo
	//send msg target_user_id
	toResp := &protopb.ApplyVsResp{
		UserInfo:       &userInfo,
		TargetUserInfo: &targetUserInfo,
		Status:         0,
		Room:           room,
	}

	// create room
	ok := this.AddRoomCache(userInfo.Id, targetUserInfo.Id, 0)
	if !ok {
		return errors.New("create room error")
	}
	ws.SendMsg(req.TargetUserId, cmd, reqUid, toResp)
	return nil
}

func (this *Wz) VsResp(req *protopb.VsRespReq, resp *protopb.VsRespResp) error {
	fmt.Println(req)
	if req.Action == 1 {
		//accept: in room
		fmt.Println("change room status:")
		ok := this.EditRoomStatusCache(req.Room, 1)
		if !ok {
			return errors.New("start vs failed")
		}
		fmt.Println("start:")
	} else {
		//reject

	}
	return nil
}

func (this *Wz) Chess(req *protopb.ChessReq, resp *protopb.ChessResp) error {
	steps := make([]Chess, 0)

	//room := this.GetRoomCache(req.Room)
	steps = append(steps, Chess{
		I: req.X,
		J: req.Y,
	})

	//check
	isWin, color := CheckFiveOfLastStep(&steps)

	resp.IsWin = isWin
	resp.Color = color
	return nil
}

func (this *Wz) Surrender(req *protopb.SurrenderReq, resp *protopb.SurrenderResp) error {

	return nil
}

func (this *Wz) DeleteUserCache() {

}

func (this *Wz) GetUserCacheById(userId uint64) (wz *types.Wz, err error) {
	key := fmt.Sprintf("%s%d", "wz:user:", userId)
	userJson := global.CACHE.Get(key)
	fmt.Println(len([]byte(userJson)))
	fmt.Println(userJson)
	fmt.Println(key)
	wz = &types.Wz{}
	if err = json.Unmarshal([]byte(userJson), wz); err != nil {
		global.LOG.Error("json反序列化失败:", err)
		return wz, err
	}
	return wz, nil
}

func (this *Wz) AddRoomCache(userId uint64, targetUserId uint64, Status uint64) bool {
	key := fmt.Sprintf("%s%d", "wz:room:", userId)
	data := types.Room{
		Id:       userId,
		TargetId: targetUserId,
		Status:   Status,
		Cap:      2,
	}
	value, err := json.Marshal(data)
	if err != nil {
		global.LOG.Error("json序列化失败: ", err)
		return false
	}
	ttl := 3600 * 24 * time.Second
	if ok := global.CACHE.Set(key, string(value), ttl); !ok {
		global.LOG.Error("create room error: ", err)
		return false
	}
	return true
}

func (this *Wz) GetRoomCache(room uint64) types.Room {
	key := fmt.Sprintf("%s%d", "wz:room:", room)
	roomRedis := global.CACHE.Get(key)
	var data types.Room
	err := json.Unmarshal([]byte(roomRedis), &data)
	if err != nil {
		global.LOG.Error("json反序列化失败: ", err)
		return data
	}
	return data
}

func (this *Wz) EditRoomStatusCache(room uint64, status uint64) bool {
	fmt.Println(room)
	key := fmt.Sprintf("%s%d", "wz:room:", room)
	data := this.GetRoomCache(room)
	fmt.Println(data)
	data.Status = status
	value, err := json.Marshal(data)
	if err != nil {
		global.LOG.Error("json序列化失败: ", err)
		return false
	}
	// delete old
	global.CACHE.Delete(key)
	// set new
	ttl := 3600 * 24 * time.Second
	if ok := global.CACHE.Set(key, string(value), ttl); !ok {
		global.LOG.Error("create room error: ", err)
		return false
	}
	return true
}

func HasStep(i int32, j int32, color int32, steps *[]Chess) bool {
	for k := int(color); k < len(*steps); k += 2 {
		step := (*steps)[k]
		if step.I == i && step.J == j {
			return true
		}
	}
	return false
}

func checkFiveInDirection(i int32, j int32, color int32, x int32, y int32, steps *[]Chess) bool {
	count := 1
	for m, n := i-x, j-y; m >= 0 && n >= 0 && m < 15 && n < 15; m, n = m-x, n-y {
		if HasStep(m, n, color, steps) {
			count++
		} else {
			break
		}
	}
	for m, n := i+x, j+y; m >= 0 && n >= 0 && m < 15 && n < 15; m, n = m+x, n+y {
		if HasStep(m, n, color, steps) {
			count++
		} else {
			break
		}
	}

	return count >= 5
}

func CheckFiveOfLastStep(steps *[]Chess) (bool, int32) {
	color := int32((len(*steps) - 1) % 2)
	if len(*steps) < 9 {
		return false, color
	}
	lastStep := (*steps)[len(*steps)-1]
	i := lastStep.I
	j := lastStep.J

	hasFive := checkFiveInDirection(i, j, color, 1, 0, steps) ||
		checkFiveInDirection(i, j, color, 0, 1, steps) ||
		checkFiveInDirection(i, j, color, 1, 1, steps) ||
		checkFiveInDirection(i, j, color, 1, -1, steps)

	return hasFive, color
}
